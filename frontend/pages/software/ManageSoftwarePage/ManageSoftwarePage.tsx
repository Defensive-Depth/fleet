import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useQuery } from "react-query";
import { InjectedRouter } from "react-router/lib/Router";
import { useDebouncedCallback } from "use-debounce";

import { AppContext } from "context/app";
import { NotificationContext } from "context/notification";
import { IConfig } from "interfaces/config";
import {
  IJiraIntegration,
  IZendeskIntegration,
  IIntegration,
} from "interfaces/integration";
import { IWebhookSoftwareVulnerabilities } from "interfaces/webhook"; // @ts-ignore
import configAPI from "services/entities/config";
import softwareAPI, {
  ISoftwareResponse,
  ISoftwareCountResponse,
} from "services/entities/software";
import {
  GITHUB_NEW_ISSUE_LINK,
  VULNERABLE_DROPDOWN_OPTIONS,
} from "utilities/constants";

import Button from "components/buttons/Button";
// @ts-ignore
import Dropdown from "components/forms/fields/Dropdown";
// @ts-ignore
import Spinner from "components/Spinner";
import TableContainer, { ITableQueryData } from "components/TableContainer";
import TableDataError from "components/DataError";
import TeamsDropdownHeader, {
  ITeamsDropdownState,
} from "components/PageHeader/TeamsDropdownHeader";
import LastUpdatedText from "components/LastUpdatedText";
import MainContent from "components/MainContent";

import generateSoftwareTableHeaders from "./SoftwareTableConfig";
import ManageAutomationsModal from "./components/ManageAutomationsModal";
import EmptySoftware from "../components/EmptySoftware";
import ExternalLinkIcon from "../../../../assets/images/open-new-tab-12x12@2x.png";

interface IManageSoftwarePageProps {
  router: InjectedRouter;
  location: {
    pathname: string;
    query: { vulnerable?: boolean };
    search: string;
  };
}

interface ISoftwareQueryKey {
  scope: string;
  page: number;
  perPage: number;
  query: string;
  orderKey: string;
  orderDir?: "asc" | "desc";
  vulnerable: boolean;
  teamId?: number;
}

interface ISoftwareAutomations {
  webhook_settings: {
    vulnerabilities_webhook: IWebhookSoftwareVulnerabilities;
  };
  integrations: {
    jira: IJiraIntegration[];
    zendesk: IZendeskIntegration[];
  };
}
interface IHeaderButtonsState extends ITeamsDropdownState {
  isLoading: boolean;
}
const DEFAULT_SORT_DIRECTION = "desc";
const DEFAULT_PAGE_SIZE = 20;

const baseClass = "manage-software-page";

const ManageSoftwarePage = ({
  router,
  location,
}: IManageSoftwarePageProps): JSX.Element => {
  const {
    availableTeams,
    currentTeam,
    isGlobalAdmin,
    isGlobalMaintainer,
    isOnGlobalTeam,
    isPremiumTier,
  } = useContext(AppContext);
  const { renderFlash } = useContext(NotificationContext);

  const DEFAULT_SORT_HEADER = isPremiumTier ? "vulnerabilities" : "hosts_count";

  const [isSoftwareEnabled, setIsSoftwareEnabled] = useState<boolean>();
  const [
    isVulnerabilityAutomationsEnabled,
    setIsVulnerabilityAutomationsEnabled,
  ] = useState<boolean>();
  const [
    recentVulnerabilityMaxAge,
    setRecentVulnerabilityMaxAge,
  ] = useState<number>();
  const [filterVuln, setFilterVuln] = useState(
    location?.query?.vulnerable || false
  );
  const [searchQuery, setSearchQuery] = useState("");
  const [sortDirection, setSortDirection] = useState<
    "asc" | "desc" | undefined
  >(DEFAULT_SORT_DIRECTION);
  const [sortHeader, setSortHeader] = useState(DEFAULT_SORT_HEADER);
  const [pageIndex, setPageIndex] = useState(0);
  const [showManageAutomationsModal, setShowManageAutomationsModal] = useState(
    false
  );
  const [showPreviewPayloadModal, setShowPreviewPayloadModal] = useState(false);

  // TODO: experiment to see if we need this state and effect or can we rely solely on the router/location for the dropdown state?
  useEffect(() => {
    setFilterVuln(!!location.query.vulnerable);
  }, [location]);

  const { data: config } = useQuery(["config"], configAPI.loadAll, {
    onSuccess: (data) => {
      setIsSoftwareEnabled(data?.features?.enable_software_inventory);
      let jiraIntegrationEnabled = false;
      if (data.integrations.jira) {
        jiraIntegrationEnabled = data?.integrations.jira.some(
          (integration: IIntegration) => {
            return integration.enable_software_vulnerabilities;
          }
        );
      }
      let zendeskIntegrationEnabled = false;
      if (data.integrations.zendesk) {
        zendeskIntegrationEnabled = data?.integrations.zendesk.some(
          (integration: IIntegration) => {
            return integration.enable_software_vulnerabilities;
          }
        );
      }
      setIsVulnerabilityAutomationsEnabled(
        data?.webhook_settings?.vulnerabilities_webhook
          .enable_vulnerabilities_webhook ||
          jiraIntegrationEnabled ||
          zendeskIntegrationEnabled
      );
      // Convert from nanosecond to nearest day
      setRecentVulnerabilityMaxAge(
        Math.round(
          data?.vulnerabilities?.recent_vulnerability_max_age / 86400000000000
        )
      );
    },
  });

  const {
    data: software,
    error: softwareError,
    isFetching: isFetchingSoftware,
  } = useQuery<
    ISoftwareResponse,
    Error,
    ISoftwareResponse,
    ISoftwareQueryKey[]
  >(
    [
      {
        scope: "software",
        page: pageIndex,
        perPage: DEFAULT_PAGE_SIZE,
        query: searchQuery,
        orderDir: sortDirection || DEFAULT_SORT_DIRECTION,
        // API expects "epss_probability" rather than "vulnerabilities"
        orderKey:
          isPremiumTier && sortHeader === "vulnerabilities"
            ? "epss_probability"
            : sortHeader,
        teamId: currentTeam?.id,
        vulnerable: !!location.query.vulnerable,
      },
    ],
    ({ queryKey }) => softwareAPI.load(queryKey[0]),
    {
      enabled:
        isOnGlobalTeam ||
        !!availableTeams?.find((t) => t.id === currentTeam?.id),
      keepPreviousData: true,
      staleTime: 30000, // stale time can be adjusted if fresher data is desired based on software inventory interval
    }
  );

  const {
    data: softwareCount,
    error: softwareCountError,
    isFetching: isFetchingCount,
  } = useQuery<
    ISoftwareCountResponse,
    Error,
    number,
    Partial<ISoftwareQueryKey>[]
  >(
    [
      {
        scope: "softwareCount",
        query: searchQuery,
        vulnerable: !!location.query.vulnerable,
        teamId: currentTeam?.id,
      },
    ],
    ({ queryKey }) => {
      return softwareAPI.count(queryKey[0]);
    },
    {
      enabled:
        isOnGlobalTeam ||
        !!availableTeams?.find((t) => t.id === currentTeam?.id),
      keepPreviousData: true,
      staleTime: 30000, // stale time can be adjusted if fresher data is desired based on software inventory interval
      refetchOnWindowFocus: false,
      retry: 1,
      select: (data) => data.count,
    }
  );

  const canAddOrRemoveSoftwareWebhook = isGlobalAdmin || isGlobalMaintainer;

  const {
    data: softwareVulnerabilitiesWebhook,
    isLoading: isLoadingSoftwareVulnerabilitiesWebhook,
    refetch: refetchSoftwareVulnerabilitiesWebhook,
  } = useQuery<IConfig, Error, IWebhookSoftwareVulnerabilities>(
    ["config"],
    () => configAPI.loadAll(),
    {
      enabled: canAddOrRemoveSoftwareWebhook,
      select: (data: IConfig) => data.webhook_settings.vulnerabilities_webhook,
    }
  );

  const onQueryChange = useDebouncedCallback(
    async ({
      pageIndex: newPageIndex,
      searchQuery: newSearchQuery,
      sortDirection: newSortDirection,
      sortHeader: newSortHeader,
    }: ITableQueryData) => {
      pageIndex !== newPageIndex && setPageIndex(newPageIndex);
      searchQuery !== newSearchQuery && setSearchQuery(newSearchQuery);
      sortDirection !== newSortDirection &&
        setSortDirection(
          newSortDirection === "asc" || newSortDirection === "desc"
            ? newSortDirection
            : DEFAULT_SORT_DIRECTION
        );

      if (isPremiumTier && newSortHeader === "vulnerabilities") {
        newSortHeader = "epss_probability";
      }
      sortHeader !== newSortHeader && setSortHeader(newSortHeader);
    },
    300
  );

  const toggleManageAutomationsModal = () =>
    setShowManageAutomationsModal(!showManageAutomationsModal);

  const togglePreviewPayloadModal = useCallback(() => {
    setShowPreviewPayloadModal(!showPreviewPayloadModal);
  }, [setShowPreviewPayloadModal, showPreviewPayloadModal]);

  const onManageAutomationsClick = () => {
    toggleManageAutomationsModal();
  };

  const onCreateWebhookSubmit = async (
    configSoftwareAutomations: ISoftwareAutomations
  ) => {
    try {
      const request = configAPI.update(configSoftwareAutomations);
      await request.then(() => {
        renderFlash(
          "success",
          "Successfully updated vulnerability automations."
        );
        refetchSoftwareVulnerabilitiesWebhook();
      });
    } catch {
      renderFlash(
        "error",
        "Could not update vulnerability automations. Please try again."
      );
    } finally {
      toggleManageAutomationsModal();
    }
  };

  const onTeamSelect = () => {
    setPageIndex(0);
  };

  const renderHeaderButtons = (
    state: IHeaderButtonsState
  ): JSX.Element | null => {
    if (
      !softwareError &&
      state.isGlobalAdmin &&
      (!state.isPremiumTier || state.teamId === 0) &&
      !state.isLoading
    ) {
      return (
        <Button
          onClick={onManageAutomationsClick}
          className={`${baseClass}__manage-automations button`}
          variant="brand"
        >
          <span>Manage automations</span>
        </Button>
      );
    }
    return null;
  };

  const renderHeaderDescription = (state: ITeamsDropdownState) => {
    return (
      <p>
        Search for installed software{" "}
        {(state.isGlobalAdmin || state.isGlobalMaintainer) &&
          (!state.isPremiumTier || state.teamId === 0) &&
          "and manage automations for detected vulnerabilities (CVEs)"}{" "}
        on{" "}
        <b>
          {state.isPremiumTier && !!state.teamId
            ? "all hosts assigned to this team"
            : "all of your hosts"}
        </b>
        .
      </p>
    );
  };

  const renderHeader = useCallback(() => {
    return (
      <TeamsDropdownHeader
        location={location}
        router={router}
        baseClass={baseClass}
        onChange={onTeamSelect}
        defaultTitle="Software"
        description={renderHeaderDescription}
        buttons={(state) =>
          renderHeaderButtons({
            ...state,
            isLoading: isLoadingSoftwareVulnerabilitiesWebhook,
          })
        }
      />
    );
  }, [router, location, isLoadingSoftwareVulnerabilitiesWebhook]);

  const renderSoftwareCount = useCallback(() => {
    const count = softwareCount;
    const lastUpdatedAt = software?.counts_updated_at;

    if (!isSoftwareEnabled || !lastUpdatedAt) {
      return null;
    }

    if (softwareCountError && !isFetchingCount) {
      return (
        <span className={`${baseClass}__count count-error`}>
          Failed to load software count
        </span>
      );
    }

    // TODO: Use setInterval to keep last updated time current?
    if (count) {
      return (
        <div
          className={`${baseClass}__count ${
            isFetchingCount ? "count-loading" : ""
          }`}
        >
          <span>{`${count} software item${count === 1 ? "" : "s"}`}</span>
          <LastUpdatedText
            lastUpdatedAt={lastUpdatedAt}
            whatToRetrieve={"software"}
          />
        </div>
      );
    }

    return null;
  }, [isFetchingCount, software, softwareCountError, softwareCount]);

  // TODO: retool this with react-router location descriptor objects
  const buildUrlQueryString = (queryString: string, vulnerable: boolean) => {
    queryString = queryString.startsWith("?")
      ? queryString.slice(1)
      : queryString;
    const queryParams = queryString.split("&").filter((el) => el.includes("="));
    const index = queryParams.findIndex((el) => el.includes("vulnerable"));

    if (vulnerable) {
      const vulnParam = `vulnerable=${vulnerable}`;
      if (index >= 0) {
        // replace old vuln param
        queryParams.splice(index, 1, vulnParam);
      } else {
        // add new vuln param
        queryParams.push(vulnParam);
      }
    } else {
      // remove old vuln param
      index >= 0 && queryParams.splice(index, 1);
    }
    queryString = queryParams.length ? "?".concat(queryParams.join("&")) : "";

    return queryString;
  };

  const onVulnFilterChange = useCallback(
    (vulnerable: boolean) => {
      setFilterVuln(vulnerable);
      setPageIndex(0);
      const queryString = buildUrlQueryString(location?.search, vulnerable);
      if (location?.search !== queryString) {
        const path = location?.pathname?.concat(queryString);
        !!path && router.replace(path);
      }
    },
    [location, router]
  );

  const renderVulnFilterDropdown = () => {
    return (
      <Dropdown
        value={filterVuln}
        className={`${baseClass}__vuln_dropdown`}
        options={VULNERABLE_DROPDOWN_OPTIONS}
        searchable={false}
        onChange={onVulnFilterChange}
      />
    );
  };

  const renderTableFooter = () => {
    return (
      <div>
        Seeing unexpected software or vulnerabilities?{" "}
        <a
          href={GITHUB_NEW_ISSUE_LINK}
          target="_blank"
          rel="noopener noreferrer"
        >
          File an issue on GitHub
          <img alt="External link" src={ExternalLinkIcon} />
        </a>
      </div>
    );
  };

  // TODO: Rework after backend is adjusted to differentiate empty search/filter results from
  // collecting inventory
  const isCollectingInventory =
    !searchQuery &&
    !filterVuln &&
    !currentTeam?.id &&
    !pageIndex &&
    !software?.software &&
    software?.counts_updated_at === null;

  const isLastPage =
    !!softwareCount &&
    DEFAULT_PAGE_SIZE * pageIndex + (software?.software?.length || 0) >=
      softwareCount;

  const softwareTableHeaders = useMemo(
    () => generateSoftwareTableHeaders(isPremiumTier),
    [isPremiumTier]
  );

  return !availableTeams || !config ? (
    <Spinner />
  ) : (
    <MainContent>
      <div className={`${baseClass}__wrapper`}>
        {renderHeader()}
        <div className={`${baseClass}__table`}>
          {softwareError && !isFetchingSoftware ? (
            <TableDataError />
          ) : (
            <TableContainer
              columns={softwareTableHeaders}
              data={(isSoftwareEnabled && software?.software) || []}
              isLoading={isFetchingSoftware || isFetchingCount}
              resultsTitle={"software items"}
              emptyComponent={() =>
                EmptySoftware(
                  (!isSoftwareEnabled && "disabled") ||
                    (isCollectingInventory && "collecting") ||
                    "default"
                )
              }
              defaultSortHeader={DEFAULT_SORT_HEADER}
              defaultSortDirection={DEFAULT_SORT_DIRECTION}
              manualSortBy
              pageSize={DEFAULT_PAGE_SIZE}
              showMarkAllPages={false}
              isAllPagesSelected={false}
              disableNextPage={isLastPage}
              searchable
              inputPlaceHolder="Search software by name or vulnerabilities (CVEs)"
              onQueryChange={onQueryChange}
              additionalQueries={filterVuln ? "vulnerable" : ""} // additionalQueries serves as a trigger
              // for the useDeepEffect hook to fire onQueryChange for events happeing outside of
              // the TableContainer
              customControl={renderVulnFilterDropdown}
              stackControls
              renderCount={renderSoftwareCount}
              renderFooter={renderTableFooter}
              disableActionButton
              hideActionButton
              highlightOnHover
            />
          )}
        </div>
        {showManageAutomationsModal && (
          <ManageAutomationsModal
            onCancel={toggleManageAutomationsModal}
            onCreateWebhookSubmit={onCreateWebhookSubmit}
            togglePreviewPayloadModal={togglePreviewPayloadModal}
            showPreviewPayloadModal={showPreviewPayloadModal}
            softwareVulnerabilityAutomationEnabled={
              isVulnerabilityAutomationsEnabled
            }
            softwareVulnerabilityWebhookEnabled={
              softwareVulnerabilitiesWebhook &&
              softwareVulnerabilitiesWebhook.enable_vulnerabilities_webhook
            }
            currentDestinationUrl={
              (softwareVulnerabilitiesWebhook &&
                softwareVulnerabilitiesWebhook.destination_url) ||
              ""
            }
            recentVulnerabilityMaxAge={recentVulnerabilityMaxAge}
          />
        )}
      </div>
    </MainContent>
  );
};

export default ManageSoftwarePage;
