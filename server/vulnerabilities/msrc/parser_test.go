package msrc

import (
	"compress/bzip2"
	"io"
	"os"
	"path/filepath"
	"testing"

	msrc_input "github.com/fleetdm/fleet/v4/server/vulnerabilities/msrc/input"
	"github.com/stretchr/testify/require"
)

func extractXMLFixtureFile(t *testing.T, src, dst string) {
	srcF, err := os.Open(src)
	require.NoError(t, err)
	defer srcF.Close()
	dstF, err := os.Create(dst)
	require.NoError(t, err)
	defer dstF.Close()
	r := bzip2.NewReader(srcF)
	// ignoring "G110: Potential DoS vulnerability via decompression bomb", as this is test code.
	_, err = io.Copy(dstF, r) //nolint:gosec
	require.NoError(t, err)
}

func TestParser(t *testing.T) {
	xmlSrcPath := filepath.Join("..", "testdata", "msrc-2022-may.xml.bz2")
	xmlDstPath := filepath.Join(t.TempDir(), "msrc-2022-may.xml")

	extractXMLFixtureFile(t, xmlSrcPath, xmlDstPath)
	f, err := os.Open(xmlDstPath)
	require.NoError(t, err)

	// Parse XML
	xmlResult, err := parseMSRCXML(f)
	f.Close()
	require.NoError(t, err)

	// All the products we expect to see, grouped by their product name
	expectedProducts := map[string]map[string]string{
		"Windows 10": {
			"11568": "Windows 10 Version 1809 for 32-bit Systems",
			"11569": "Windows 10 Version 1809 for x64-based Systems",
			"11570": "Windows 10 Version 1809 for ARM64-based Systems",
			"11712": "Windows 10 Version 1909 for 32-bit Systems",
			"11713": "Windows 10 Version 1909 for x64-based Systems",
			"11714": "Windows 10 Version 1909 for ARM64-based Systems",
			"11896": "Windows 10 Version 21H1 for x64-based Systems",
			"11897": "Windows 10 Version 21H1 for ARM64-based Systems",
			"11898": "Windows 10 Version 21H1 for 32-bit Systems",
			"11800": "Windows 10 Version 20H2 for x64-based Systems",
			"11801": "Windows 10 Version 20H2 for 32-bit Systems",
			"11802": "Windows 10 Version 20H2 for ARM64-based Systems",
			"11929": "Windows 10 Version 21H2 for 32-bit Systems",
			"11930": "Windows 10 Version 21H2 for ARM64-based Systems",
			"11931": "Windows 10 Version 21H2 for x64-based Systems",
			"10729": "Windows 10 for 32-bit Systems",
			"10735": "Windows 10 for x64-based Systems",
			"10852": "Windows 10 Version 1607 for 32-bit Systems",
			"10853": "Windows 10 Version 1607 for x64-based Systems",
		},
		"Windows Server 2019": {
			"11571": "Windows Server 2019",
			"11572": "Windows Server 2019  (Server Core installation)",
		},
		"Windows Server 2022": {
			"11923": "Windows Server 2022",
			"11924": "Windows Server 2022 (Server Core installation)",
		},
		"Windows Server": {
			"11803": "Windows Server, version 20H2 (Server Core Installation)",
		},
		"Windows 11": {
			"11926": "Windows 11 for x64-based Systems",
			"11927": "Windows 11 for ARM64-based Systems",
		},
		"Windows Server 2016": {
			"10816": "Windows Server 2016",
			"10855": "Windows Server 2016  (Server Core installation)",
		},
		"Windows 8.1": {
			"10481": "Windows 8.1 for 32-bit systems",
			"10482": "Windows 8.1 for x64-based systems",
		},
		"Windows RT 8.1": {
			"10484": "Windows RT 8.1",
		},
		"Windows Server 2012": {
			"10378": "Windows Server 2012",
			"10379": "Windows Server 2012 (Server Core installation)",
		},
		"Windows Server 2012 R2": {
			"10483": "Windows Server 2012 R2",
			"10543": "Windows Server 2012 R2 (Server Core installation)",
		},
		"Windows 7": {
			"10047": "Windows 7 for 32-bit Systems Service Pack 1",
			"10048": "Windows 7 for x64-based Systems Service Pack 1",
		},
		"Windows Server 2008": {
			"9312":  "Windows Server 2008 for 32-bit Systems Service Pack 2",
			"10287": "Windows Server 2008 for 32-bit Systems Service Pack 2 (Server Core installation)",
			"9318":  "Windows Server 2008 for x64-based Systems Service Pack 2",
			"9344":  "Windows Server 2008 for x64-based Systems Service Pack 2 (Server Core installation)",
		},
		"Windows Server 2008 R2": {
			"10051": "Windows Server 2008 R2 for x64-based Systems Service Pack 1",
			"10049": "Windows Server 2008 R2 for x64-based Systems Service Pack 1 (Server Core installation)",
		},
	}

	expectedCVEs := map[string][]string{
		"Windows 10": {
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-23279",
			"CVE-2022-29142",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-29140",
			"CVE-2022-21972",
			"CVE-2022-22713",
			"CVE-2022-23270",
			"CVE-2022-24466",
			"CVE-2022-26913",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26927",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-22016",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29105",
			"CVE-2022-29112",
			"CVE-2022-29113",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29131",
			"CVE-2022-29132",
			"CVE-2022-29137",
			"CVE-2022-29139",
		},
		"Windows Server 2019": {
			"CVE-2022-26927",
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-29142",
			"CVE-2022-29150",
			"CVE-2022-29151",
			"CVE-2022-29122",
			"CVE-2022-29120",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-29140",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-24466",
			"CVE-2022-26913",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26932",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-26937",
			"CVE-2022-26938",
			"CVE-2022-26939",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-22016",
			"CVE-2022-29102",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29105",
			"CVE-2022-29106",
			"CVE-2022-29112",
			"CVE-2022-29113",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29131",
			"CVE-2022-29132",
			"CVE-2022-29134",
			"CVE-2022-29135",
			"CVE-2022-29137",
			"CVE-2022-29138",
			"CVE-2022-29123",
			"CVE-2022-29139",
		},
		"Windows Server 2022": {
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-23279",
			"CVE-2022-29142",
			"CVE-2022-29150",
			"CVE-2022-29151",
			"CVE-2022-29122",
			"CVE-2022-29120",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-29140",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-24466",
			"CVE-2022-26913",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26927",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26932",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-26937",
			"CVE-2022-26938",
			"CVE-2022-26939",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-22016",
			"CVE-2022-29102",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29106",
			"CVE-2022-29112",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29131",
			"CVE-2022-29132",
			"CVE-2022-29134",
			"CVE-2022-29135",
			"CVE-2022-29137",
			"CVE-2022-29138",
			"CVE-2022-29123",
			"CVE-2022-22017",
			"CVE-2022-26940",
			"CVE-2022-29139",
		},
		"Windows Server": {
			"CVE-2022-24466",
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-23279",
			"CVE-2022-29142",
			"CVE-2022-29150",
			"CVE-2022-29151",
			"CVE-2022-29122",
			"CVE-2022-29120",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-29140",
			"CVE-2022-21972",
			"CVE-2022-22713",
			"CVE-2022-23270",
			"CVE-2022-26913",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26927",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26932",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-26937",
			"CVE-2022-26938",
			"CVE-2022-26939",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-22016",
			"CVE-2022-29102",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29105",
			"CVE-2022-29106",
			"CVE-2022-29112",
			"CVE-2022-29113",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29131",
			"CVE-2022-29132",
			"CVE-2022-29134",
			"CVE-2022-29135",
			"CVE-2022-29137",
			"CVE-2022-29138",
			"CVE-2022-29123",
			"CVE-2022-29139",
		},
		"Windows 11": {
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-23279",
			"CVE-2022-29116",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-29140",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-24466",
			"CVE-2022-26913",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26927",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-22016",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29112",
			"CVE-2022-29113",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29131",
			"CVE-2022-29132",
			"CVE-2022-29133",
			"CVE-2022-29137",
			"CVE-2022-22017",
			"CVE-2022-26940",
			"CVE-2022-29139",
		},
		"Windows Server 2016": {
			"CVE-2022-29137",
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-29150",
			"CVE-2022-29151",
			"CVE-2022-29122",
			"CVE-2022-29120",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-24466",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26932",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-26937",
			"CVE-2022-26938",
			"CVE-2022-26939",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-22016",
			"CVE-2022-29102",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29105",
			"CVE-2022-29106",
			"CVE-2022-29112",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29132",
			"CVE-2022-29134",
			"CVE-2022-29135",
			"CVE-2022-29138",
			"CVE-2022-29123",
			"CVE-2022-29139",
			"CVE-2022-29140",
		},
		"Windows 8.1": {
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29105",
			"CVE-2022-29112",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29132",
			"CVE-2022-29137",
			"CVE-2022-29139",
		},
		"Windows RT 8.1": {
			"CVE-2022-26934",
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26933",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29105",
			"CVE-2022-29112",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29132",
			"CVE-2022-29137",
			"CVE-2022-29139",
		},
		"Windows Server 2012": {
			"CVE-2022-26936",
			"CVE-2022-30190",
			"CVE-2022-29150",
			"CVE-2022-29151",
			"CVE-2022-29122",
			"CVE-2022-29120",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26937",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-29102",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29105",
			"CVE-2022-29112",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29132",
			"CVE-2022-29135",
			"CVE-2022-29137",
			"CVE-2022-29138",
			"CVE-2022-29123",
			"CVE-2022-29139",
		},
		"Windows Server 2012 R2": {
			"CVE-2022-30190",
			"CVE-2022-26923",
			"CVE-2022-29150",
			"CVE-2022-29151",
			"CVE-2022-29122",
			"CVE-2022-29120",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26930",
			"CVE-2022-26931",
			"CVE-2022-26933",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26937",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-29102",
			"CVE-2022-29103",
			"CVE-2022-29104",
			"CVE-2022-29105",
			"CVE-2022-29112",
			"CVE-2022-29114",
			"CVE-2022-29115",
			"CVE-2022-29125",
			"CVE-2022-29126",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29132",
			"CVE-2022-29134",
			"CVE-2022-29135",
			"CVE-2022-29137",
			"CVE-2022-29138",
			"CVE-2022-29123",
			"CVE-2022-29139",
			"CVE-2022-26936",
		},
		"Windows 7": {
			"CVE-2022-29105",
			"CVE-2022-30190",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26931",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-29103",
			"CVE-2022-29112",
			"CVE-2022-29115",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29132",
			"CVE-2022-29137",
			"CVE-2022-29139",
		},
		"Windows Server 2008": {
			"CVE-2022-29115",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26931",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-26937",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-29103",
			"CVE-2022-29112",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29132",
			"CVE-2022-29137",
			"CVE-2022-29139",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
		},
		"Windows Server 2008 R2": {
			"CVE-2022-30190",
			"CVE-2022-21972",
			"CVE-2022-23270",
			"CVE-2022-26925",
			"CVE-2022-26926",
			"CVE-2022-26931",
			"CVE-2022-26934",
			"CVE-2022-26935",
			"CVE-2022-26936",
			"CVE-2022-26937",
			"CVE-2022-22011",
			"CVE-2022-22012",
			"CVE-2022-22013",
			"CVE-2022-22014",
			"CVE-2022-22015",
			"CVE-2022-29103",
			"CVE-2022-29112",
			"CVE-2022-29115",
			"CVE-2022-29127",
			"CVE-2022-29128",
			"CVE-2022-29129",
			"CVE-2022-29130",
			"CVE-2022-29132",
			"CVE-2022-29137",
			"CVE-2022-29139",
			"CVE-2022-29141",
			"CVE-2022-22019",
			"CVE-2022-29121",
			"CVE-2022-30138",
			"CVE-2022-29105",
		},
	}

	t.Run("mapToVulnGraphs", func(t *testing.T) {
		graphs, err := mapToVulnGraphs(xmlResult)
		require.NoError(t, err)

		t.Run("should have one graph per product", func(t *testing.T) {
			var expected []string
			for p := range expectedProducts {
				expected = append(expected, p)
			}

			var actual []string
			for _, g := range graphs {
				actual = append(actual, g.ProductName)
			}

			require.Len(t, graphs, len(expected))
			require.ElementsMatch(t, expected, actual)
		})

		t.Run("each graph should have the right products", func(t *testing.T) {
			for _, g := range graphs {
				require.Equal(t, g.Products, expectedProducts[g.ProductName], g.ProductName)
			}
		})

		t.Run("each graph should have the right vulnerabilities", func(t *testing.T) {
			for _, g := range graphs {
				var actual []string
				for v := range g.Vulnerabities {
					actual = append(actual, v)
				}
				require.ElementsMatch(t, actual, expectedCVEs[g.ProductName], g.ProductName)
			}
		})
	})

	t.Run("parseMSRCXML", func(t *testing.T) {
		t.Run("only windows products are included", func(t *testing.T) {
			var expected []msrc_input.ProductXML
			for _, grp := range expectedProducts {
				for pID, pFn := range grp {
					expected = append(
						expected,
						msrc_input.ProductXML{ProductID: pID, FullName: pFn},
					)
				}
			}

			var actual []msrc_input.ProductXML
			for _, v := range xmlResult.WinProducts {
				actual = append(actual, v)
			}
			require.ElementsMatch(t, actual, expected)
		})

		t.Run("only CVEs for windows products are included", func(t *testing.T) {
			expected := make(map[string]bool)
			for _, p := range expectedCVEs {
				for _, v := range p {
					expected[v] = true
				}
			}
			actual := make(map[string]bool)
			for _, v := range xmlResult.WinVulnerabities {
				actual[v.CVE] = true
			}
			require.Equal(t, expected, actual)
		})

		t.Run("scores are parsed correctly", func(t *testing.T) {
			// Check the score of a random CVE (CVE-2022-24466)
			for _, v := range xmlResult.WinVulnerabities {
				if v.CVE == "CVE-2022-24466" {
					require.Equal(t, 4.1, v.Score)
				}
			}
		})

		t.Run("the revision history is parsed correctly", func(t *testing.T) {
			// Check the revision history of a random CVE (CVE-2022-29114)
			for _, v := range xmlResult.WinVulnerabities {
				if v.CVE == "CVE-2022-29114" {
					require.Len(t, v.Revisions, 1)
					require.Equal(t, "2022-05-10T08:00:00", v.Revisions[0].Date)
					require.Equal(t, "<p>Information published.</p>\n", v.Revisions[0].Description)
				}
			}
		})

		t.Run("the remediations are parsed correctly", func(t *testing.T) {
			// Check the remediations of a random CVE (CVE-2022-29126)
			expectedRemediations := []msrc_input.VulnerabilityRemediationXML{
				{Type: "Vendor Fix", FixedBuild: "10.0.17763.2928", ProductIDs: []string{"11568", "11569", "11570", "11571", "11572"}, Description: "5013941", Supercedence: "5012647", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"11568", "11569", "11570", "11571", "11572"}, Description: "5013941"},
				{Type: "Vendor Fix", FixedBuild: "10.0.18363.2274", ProductIDs: []string{"11712", "11713", "11714"}, Description: "5013945", Supercedence: "5012591", RestartRequired: "Yes"},
				{Type: "Vendor Fix", FixedBuild: "10.0.19043.1706", ProductIDs: []string{"11896", "11897", "11898", "11929"}, Description: "5013942", Supercedence: "5012599", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"11896", "11897", "11898", "11800", "11801", "11802", "11803", "11929", "11930", "11931"}, Description: "5013942"},
				{Type: "Vendor Fix", FixedBuild: "10.0.20348.707", ProductIDs: []string{"11923", "11924"}, Description: "5013944", Supercedence: "5012604", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"11923", "11924"}, Description: "5013944"},
				{Type: "Vendor Fix", FixedBuild: "10.0.19042.1706", ProductIDs: []string{"11800", "11801", "11802", "11803"}, Description: "5013942", Supercedence: "5012599", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"11896", "11897", "11898", "11800", "11801", "11802", "11803", "11929", "11930", "11931"}, Description: "5013942"},
				{Type: "Vendor Fix", FixedBuild: "10.0.22000.675", ProductIDs: []string{"11926", "11927"}, Description: "5013943", Supercedence: "5012592", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"11926", "11927"}, Description: "5013943"},
				{Type: "Vendor Fix", FixedBuild: "10.0.19044.1706", ProductIDs: []string{"11930", "11931"}, Description: "5013942", Supercedence: "5012599", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"11896", "11897", "11898", "11800", "11801", "11802", "11803", "11929", "11930", "11931"}, Description: "5013942"},
				{Type: "Vendor Fix", FixedBuild: "10.0.10240.19297", ProductIDs: []string{"10729", "10735"}, Description: "5013963", Supercedence: "5012653", RestartRequired: "Yes"},
				{Type: "Vendor Fix", FixedBuild: "10.0.14393.5125", ProductIDs: []string{"10852", "10853", "10816", "10855"}, Description: "5013952", Supercedence: "5012596", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"10852", "10853", "10816", "10855"}, Description: "5013952"},
				{Type: "Vendor Fix", FixedBuild: "6.3.9600.20371", ProductIDs: []string{"10481", "10482", "10483", "10543"}, Description: "5014011", Supercedence: "5012670", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"10481", "10482", "10483", "10543"}, Description: "5014011"},
				{Type: "Vendor Fix", FixedBuild: "6.3.9600.20365", ProductIDs: []string{"10481", "10482", "10483", "10543"}, Description: "5014001", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"10481", "10482", "10483", "10543"}, Description: "5014001"},
				{Type: "Vendor Fix", FixedBuild: "6.3.9600.20367", ProductIDs: []string{"10484"}, Description: "5014025", RestartRequired: "Yes"},
				{Type: "Vendor Fix", FixedBuild: "6.2.9200.23714", ProductIDs: []string{"10378", "10379"}, Description: "5014017", Supercedence: "5012650", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"10378", "10379"}, Description: "5014017"},
				{Type: "Vendor Fix", FixedBuild: "6.2.9200.23714", ProductIDs: []string{"10378", "10379"}, Description: "5014018", RestartRequired: "Yes"},
				{Type: "Known Issue", ProductIDs: []string{"10378", "10379"}, Description: "5014018"},
			}
			for _, v := range xmlResult.WinVulnerabities {
				if v.CVE == "CVE-2022-29126" {
					require.Len(t, v.Remediations, len(expectedRemediations))
					require.ElementsMatch(t, v.Remediations, expectedRemediations)
				}
			}
		})
	})
}
