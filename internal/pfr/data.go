package pfr

// Pro Football Reference team keys
var PFR_TEAM_KEYS = map[string]string{
	"ARI": "crd",
	"ATL": "atl",
	"BAL": "rav",
	"BUF": "buf",
	"CAR": "car",
	"CHI": "chi",
	"CIN": "cin",
	"CLE": "cle",
	"DAL": "dal",
	"DEN": "den",
	"DET": "det",
	"GB":  "gnb",
	"HOU": "htx",
	"IND": "clt",
	"JAX": "jax",
	"KC":  "kan",
	"LAC": "sdg",
	"LAR": "ram",
	"LVR": "rai",
	"MIA": "mia",
	"MIN": "min",
	"NO":  "nor",
	"NE":  "nwe",
	"NYG": "nyg",
	"NYJ": "nyj",
	"PHI": "phi",
	"PIT": "pit",
	"SEA": "sea",
	"SF":  "sfo",
	"TB":  "tam",
	"TEN": "oti",
	"WSH": "was",
}

// Pro Football Reference table ids
var PFR_TABLE_IDS = []string{
	"passing",
	"rushing_and_receiving",
}

// Pro Football Reference table headers to rename
var HEADER_RENAMES = map[string]string{
	"name_display":    "player",
	"games":           "g",
	"games_started":   "gs",
	"rush_first_down": "rush_1d",
	"rush_success":    "rush_1d%",
	"rec_first_down":  "rec_1d",
	"rec_success":     "rec_1d%",
}

var OUT_HEADERS = []string{
	"player",
	"age",
	"pos",
	"g",
	"gs",
	"pass_cmp",
	"pass_att",
	"pass_yds",
	"pass_td",
	"pass_int",
	"pass_sacked_yds",
	"rush_att",
	"rush_yds",
	"rush_td",
	"rush_1d",
	"rush_1d%",
	"targets",
	"rec",
	"rec_yds",
	"rec_td",
	"rec_1d",
	"rec_1d%",
	"touches",
	"fumbles",
}
