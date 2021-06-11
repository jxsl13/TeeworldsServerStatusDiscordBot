package markdown

import "strings"

var (
	flags = map[int]string{
		737: "SS",
		901: "XEN",
		902: "XNI",
		903: "XSC",
		904: "XWA",
		950: "XBZ",
		951: "XCA",
		952: "XES",
		953: "XGA",
		-1:  "default",
		20:  "AD",
		784: "AE",
		4:   "AF",
		28:  "AG",
		660: "AI",
		8:   "AL",
		51:  "AM",
		24:  "AO",
		32:  "AR",
		16:  "AS",
		40:  "AT",
		36:  "AU",
		533: "AW",
		248: "AX",
		31:  "AZ",
		70:  "BA",
		52:  "BB",
		50:  "BD",
		56:  "BE",
		854: "BF",
		100: "BG",
		48:  "BH",
		108: "BI",
		204: "BJ",
		652: "BL",
		60:  "BM",
		96:  "BN",
		68:  "BO",
		76:  "BR",
		44:  "BS",
		64:  "BT",
		72:  "BW",
		112: "BY",
		84:  "BZ",
		124: "CA",
		166: "CC",
		180: "CD",
		140: "CF",
		178: "CG",
		756: "CH",
		384: "CI",
		184: "CK",
		152: "CL",
		120: "CM",
		156: "CN",
		170: "CO",
		188: "CR",
		192: "CU",
		132: "CV",
		531: "CW",
		162: "CX",
		196: "CY",
		203: "CZ",
		276: "DE",
		262: "DJ",
		208: "DK",
		212: "DM",
		214: "DO",
		12:  "DZ",
		218: "EC",
		233: "EE",
		818: "EG",
		732: "EH",
		232: "ER",
		724: "ES",
		231: "ET",
		246: "FI",
		242: "FJ",
		238: "FK",
		583: "FM",
		234: "FO",
		250: "FR",
		266: "GA",
		826: "GB",
		308: "GD",
		268: "GE",
		254: "GF",
		831: "GG",
		288: "GH",
		292: "GI",
		304: "GL",
		270: "GM",
		324: "GN",
		312: "GP",
		226: "GQ",
		300: "GR",
		239: "GS",
		320: "GT",
		316: "GU",
		624: "GW",
		328: "GY",
		344: "HK",
		340: "HN",
		191: "HR",
		332: "HT",
		348: "HU",
		360: "ID",
		372: "IE",
		376: "IL",
		833: "IM",
		356: "IN",
		86:  "IO",
		368: "IQ",
		364: "IR",
		352: "IS",
		380: "IT",
		832: "JE",
		388: "JM",
		400: "JO",
		392: "JP",
		404: "KE",
		417: "KG",
		116: "KH",
		296: "KI",
		174: "KM",
		659: "KN",
		408: "KP",
		410: "KR",
		414: "KW",
		136: "KY",
		398: "KZ",
		418: "LA",
		422: "LB",
		662: "LC",
		438: "LI",
		144: "LK",
		430: "LR",
		426: "LS",
		440: "LT",
		442: "LU",
		428: "LV",
		434: "LY",
		504: "MA",
		492: "MC",
		498: "MD",
		499: "ME",
		663: "MF",
		450: "MG",
		584: "MH",
		807: "MK",
		466: "ML",
		104: "MM",
		496: "MN",
		446: "MO",
		580: "MP",
		474: "MQ",
		478: "MR",
		500: "MS",
		470: "MT",
		480: "MU",
		462: "MV",
		454: "MW",
		484: "MX",
		458: "MY",
		508: "MZ",
		516: "NA",
		540: "NC",
		562: "NE",
		574: "NF",
		566: "NG",
		558: "NI",
		528: "NL",
		578: "NO",
		524: "NP",
		520: "NR",
		570: "NU",
		554: "NZ",
		512: "OM",
		591: "PA",
		604: "PE",
		258: "PF",
		598: "PG",
		608: "PH",
		586: "PK",
		616: "PL",
		666: "PM",
		612: "PN",
		630: "PR",
		275: "PS",
		620: "PT",
		585: "PW",
		600: "PY",
		634: "QA",
		638: "RE",
		642: "RO",
		688: "RS",
		643: "RU",
		646: "RW",
		682: "SA",
		90:  "SB",
		690: "SC",
		736: "SD",
		752: "SE",
		702: "SG",
		654: "SH",
		705: "SI",
		703: "SK",
		694: "SL",
		674: "SM",
		686: "SN",
		706: "SO",
		740: "SR",
		678: "ST",
		222: "SV",
		534: "SX",
		760: "SY",
		748: "SZ",
		796: "TC",
		148: "TD",
		260: "TF",
		768: "TG",
		764: "TH",
		762: "TJ",
		772: "TK",
		626: "TL",
		795: "TM",
		788: "TN",
		776: "TO",
		792: "TR",
		780: "TT",
		798: "TV",
		158: "TW",
		834: "TZ",
		804: "UA",
		800: "UG",
		840: "US",
		858: "UY",
		860: "UZ",
		336: "VA",
		670: "VC",
		862: "VE",
		92:  "VG",
		850: "VI",
		704: "VN",
		548: "VU",
		876: "WF",
		882: "WS",
		887: "YE",
		710: "ZA",
		894: "ZM",
		716: "ZW",
		10:  "AQ",
		535: "BQ",
		74:  "BV",
		334: "HM",
		744: "SJ",
		581: "UM",
		175: "YT",
	}
)

// Flag returns a string representation of a given flag value
func Flag(value int) string {
	if value <= 0 {
		return ":rainbow_flag:"
	}
	return ":flag_" + strings.ToLower(flags[value]) + ":"
}
