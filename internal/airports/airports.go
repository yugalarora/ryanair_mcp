package airports

import (
	"math"
	"sort"
	"strings"
)

// Airport represents a Ryanair destination with geographical coordinates.
type Airport struct {
	IATA      string  `json:"iata"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Aliases   []string `json:"aliases,omitempty"`
}

// NearbyAirport pairs an airport with its distance from a reference location.
type NearbyAirport struct {
	Airport
	DistanceKm float64 `json:"distance_km"`
}

// Directory is the curated database of Ryanair European and North African airports.
var Directory = map[string]Airport{
	// Ireland
	"DUB": {IATA: "DUB", Name: "Dublin Airport", City: "Dublin", Country: "Ireland", Latitude: 53.4213, Longitude: -6.2701, Aliases: []string{"dublin", "leinster"}},
	"SNN": {IATA: "SNN", Name: "Shannon Airport", City: "Shannon", Country: "Ireland", Latitude: 52.7020, Longitude: -8.9248, Aliases: []string{"limerick", "galway", "shannon"}},
	"ORK": {IATA: "ORK", Name: "Cork Airport", City: "Cork", Country: "Ireland", Latitude: 51.8413, Longitude: -8.4911, Aliases: []string{"cork", "munster"}},
	"NOC": {IATA: "NOC", Name: "Ireland West Airport Knock", City: "Knock", Country: "Ireland", Latitude: 53.9103, Longitude: -8.8185, Aliases: []string{"knock", "mayo", "connacht"}},
	"BFS": {IATA: "BFS", Name: "Belfast International Airport", City: "Belfast", Country: "United Kingdom", Latitude: 54.6575, Longitude: -6.2158, Aliases: []string{"belfast", "northern ireland"}},

	// United Kingdom
	"STN": {IATA: "STN", Name: "London Stansted Airport", City: "London", Country: "United Kingdom", Latitude: 51.8860, Longitude: 0.2389, Aliases: []string{"london", "stansted", "essex"}},
	"LGW": {IATA: "LGW", Name: "London Gatwick Airport", City: "London", Country: "United Kingdom", Latitude: 51.1537, Longitude: -0.1821, Aliases: []string{"london", "gatwick", "surrey"}},
	"LTN": {IATA: "LTN", Name: "London Luton Airport", City: "London", Country: "United Kingdom", Latitude: 51.8747, Longitude: -0.3683, Aliases: []string{"london", "luton", "bedfordshire"}},
	"MAN": {IATA: "MAN", Name: "Manchester Airport", City: "Manchester", Country: "United Kingdom", Latitude: 53.3537, Longitude: -2.2750, Aliases: []string{"manchester", "north west"}},
	"BHX": {IATA: "BHX", Name: "Birmingham Airport", City: "Birmingham", Country: "United Kingdom", Latitude: 52.4524, Longitude: -1.7435, Aliases: []string{"birmingham", "midlands"}},
	"EDI": {IATA: "EDI", Name: "Edinburgh Airport", City: "Edinburgh", Country: "United Kingdom", Latitude: 55.9508, Longitude: -3.3725, Aliases: []string{"edinburgh", "scotland"}},
	"GLA": {IATA: "GLA", Name: "Glasgow Airport", City: "Glasgow", Country: "United Kingdom", Latitude: 55.8719, Longitude: -4.4331, Aliases: []string{"glasgow", "scotland"}},
	"BRS": {IATA: "BRS", Name: "Bristol Airport", City: "Bristol", Country: "United Kingdom", Latitude: 51.3827, Longitude: -2.7191, Aliases: []string{"bristol", "bath", "somerset"}},
	"LPL": {IATA: "LPL", Name: "Liverpool John Lennon Airport", City: "Liverpool", Country: "United Kingdom", Latitude: 53.3336, Longitude: -2.8497, Aliases: []string{"liverpool", "merseyside"}},
	"EMA": {IATA: "EMA", Name: "East Midlands Airport", City: "Nottingham", Country: "United Kingdom", Latitude: 52.8311, Longitude: -1.3281, Aliases: []string{"nottingham", "derby", "leicester"}},
	"CWL": {IATA: "CWL", Name: "Cardiff Airport", City: "Cardiff", Country: "United Kingdom", Latitude: 51.3967, Longitude: -3.3433, Aliases: []string{"cardiff", "wales"}},
	"NCL": {IATA: "NCL", Name: "Newcastle Airport", City: "Newcastle", Country: "United Kingdom", Latitude: 55.0375, Longitude: -1.6917, Aliases: []string{"newcastle", "tyne"}},
	"LBA": {IATA: "LBA", Name: "Leeds Bradford Airport", City: "Leeds", Country: "United Kingdom", Latitude: 53.8659, Longitude: -1.6606, Aliases: []string{"leeds", "bradford", "yorkshire"}},

	// Belgium & Luxembourg & Netherlands
	"CRL": {IATA: "CRL", Name: "Brussels South Charleroi Airport", City: "Brussels", Country: "Belgium", Latitude: 50.4592, Longitude: 4.4538, Aliases: []string{"brussels", "charleroi", "wallonia", "belgium"}},
	"BRU": {IATA: "BRU", Name: "Brussels Airport (Zaventem)", City: "Brussels", Country: "Belgium", Latitude: 50.9014, Longitude: 4.4844, Aliases: []string{"brussels", "zaventem", "belgium"}},
	"LGG": {IATA: "LGG", Name: "Liège Airport", City: "Liège", Country: "Belgium", Latitude: 50.6374, Longitude: 5.4432, Aliases: []string{"liege", "luik", "belgium"}},
	"LUX": {IATA: "LUX", Name: "Luxembourg Airport", City: "Luxembourg", Country: "Luxembourg", Latitude: 49.6233, Longitude: 6.2044, Aliases: []string{"luxembourg", "lux"}},
	"EIN": {IATA: "EIN", Name: "Eindhoven Airport", City: "Eindhoven", Country: "Netherlands", Latitude: 51.4501, Longitude: 5.3745, Aliases: []string{"eindhoven", "brabant", "netherlands"}},
	"AMS": {IATA: "AMS", Name: "Amsterdam Airport Schiphol", City: "Amsterdam", Country: "Netherlands", Latitude: 52.3105, Longitude: 4.7683, Aliases: []string{"amsterdam", "schiphol", "holland"}},
	"MST": {IATA: "MST", Name: "Maastricht Aachen Airport", City: "Maastricht", Country: "Netherlands", Latitude: 50.9117, Longitude: 5.7701, Aliases: []string{"maastricht", "aachen", "limburg"}},
	"LIL": {IATA: "LIL", Name: "Lille Airport", City: "Lille", Country: "France", Latitude: 50.5619, Longitude: 3.0894, Aliases: []string{"lille", "flanders"}},

	// Portugal
	"FAO": {IATA: "FAO", Name: "Faro Airport", City: "Faro", Country: "Portugal", Latitude: 37.0176, Longitude: -7.9659, Aliases: []string{"faro", "algarve", "portugal", "albufeira", "lagos", "vilamoura"}},
	"LIS": {IATA: "LIS", Name: "Humberto Delgado Airport (Lisbon)", City: "Lisbon", Country: "Portugal", Latitude: 38.7742, Longitude: -9.1342, Aliases: []string{"lisbon", "lisboa", "sintra", "cascais"}},
	"OPO": {IATA: "OPO", Name: "Francisco Sá Carneiro Airport (Porto)", City: "Porto", Country: "Portugal", Latitude: 41.2481, Longitude: -8.6814, Aliases: []string{"porto", "oporto", "douro"}},
	"FNC": {IATA: "FNC", Name: "Cristiano Ronaldo International Airport (Madeira)", City: "Funchal", Country: "Portugal", Latitude: 32.6979, Longitude: -16.7745, Aliases: []string{"madeira", "funchal"}},
	"PDL": {IATA: "PDL", Name: "João Paulo II Airport (Azores)", City: "Ponta Delgada", Country: "Portugal", Latitude: 37.7412, Longitude: -25.6979, Aliases: []string{"azores", "ponta delgada", "sao miguel"}},

	// Spain
	"MAD": {IATA: "MAD", Name: "Adolfo Suárez Madrid–Barajas Airport", City: "Madrid", Country: "Spain", Latitude: 40.4839, Longitude: -3.5680, Aliases: []string{"madrid", "barajas", "spain"}},
	"BCN": {IATA: "BCN", Name: "Josep Tarradellas Barcelona-El Prat Airport", City: "Barcelona", Country: "Spain", Latitude: 41.2974, Longitude: 2.0833, Aliases: []string{"barcelona", "el prat", "catalonia"}},
	"GRO": {IATA: "GRO", Name: "Girona-Costa Brava Airport", City: "Girona", Country: "Spain", Latitude: 41.9010, Longitude: 2.7606, Aliases: []string{"girona", "costa brava", "barcelona north"}},
	"REU": {IATA: "REU", Name: "Reus Airport", City: "Reus", Country: "Spain", Latitude: 41.1474, Longitude: 1.1672, Aliases: []string{"reus", "tarragona", "salou", "costa daurada"}},
	"AGP": {IATA: "AGP", Name: "Málaga-Costa del Sol Airport", City: "Málaga", Country: "Spain", Latitude: 36.6749, Longitude: -4.4991, Aliases: []string{"malaga", "costa del sol", "marbella", "torremolinos", "andalucia", "andalusia", "tarifa", "gibraltar"}},
	"SVQ": {IATA: "SVQ", Name: "Seville Airport", City: "Seville", Country: "Spain", Latitude: 37.4180, Longitude: -5.8931, Aliases: []string{"seville", "sevilla", "andalucia", "andalusia"}},
	"XRY": {IATA: "XRY", Name: "Jerez Airport", City: "Jerez de la Frontera", Country: "Spain", Latitude: 36.7447, Longitude: -6.0601, Aliases: []string{"jerez", "cadiz", "tarifa", "gibraltar", "andalucia", "andalusia"}},
	"ALC": {IATA: "ALC", Name: "Alicante–Elche Miguel Hernández Airport", City: "Alicante", Country: "Spain", Latitude: 38.2822, Longitude: -0.5582, Aliases: []string{"alicante", "costa blanca", "benidorm", "altea", "calpe", "elche"}},
	"RMU": {IATA: "RMU", Name: "Región de Murcia International Airport", City: "Murcia", Country: "Spain", Latitude: 37.8030, Longitude: -1.1250, Aliases: []string{"murcia", "corvera", "costa calida"}},
	"VLC": {IATA: "VLC", Name: "Valencia Airport", City: "Valencia", Country: "Spain", Latitude: 39.4893, Longitude: -0.4816, Aliases: []string{"valencia"}},
	"ZAZ": {IATA: "ZAZ", Name: "Zaragoza Airport", City: "Zaragoza", Country: "Spain", Latitude: 41.6662, Longitude: -1.0416, Aliases: []string{"zaragoza", "aragon"}},
	"PMI": {IATA: "PMI", Name: "Palma de Mallorca Airport", City: "Palma de Mallorca", Country: "Spain", Latitude: 39.5517, Longitude: 2.7388, Aliases: []string{"mallorca", "majorca", "palma", "balearics"}},
	"IBZ": {IATA: "IBZ", Name: "Ibiza Airport", City: "Ibiza", Country: "Spain", Latitude: 38.8729, Longitude: 1.3731, Aliases: []string{"ibiza", "eivissa", "balearics"}},
	"MAH": {IATA: "MAH", Name: "Menorca Airport", City: "Menorca", Country: "Spain", Latitude: 39.8626, Longitude: 4.2186, Aliases: []string{"menorca", "minorca", "mahon"}},
	"TFS": {IATA: "TFS", Name: "Tenerife South Airport", City: "Tenerife", Country: "Spain", Latitude: 28.0445, Longitude: -16.5725, Aliases: []string{"tenerife", "canaries", "reina sofia"}},
	"TFN": {IATA: "TFN", Name: "Tenerife North Airport", City: "Tenerife", Country: "Spain", Latitude: 28.4827, Longitude: -16.3415, Aliases: []string{"tenerife north", "los rodeos"}},
	"LPA": {IATA: "LPA", Name: "Gran Canaria Airport", City: "Gran Canaria", Country: "Spain", Latitude: 27.9319, Longitude: -15.3866, Aliases: []string{"gran canaria", "las palmas", "canaries"}},
	"ACE": {IATA: "ACE", Name: "Lanzarote Airport", City: "Lanzarote", Country: "Spain", Latitude: 28.9455, Longitude: -13.6052, Aliases: []string{"lanzarote", "arrecife", "canaries"}},
	"FUE": {IATA: "FUE", Name: "Fuerteventura Airport", City: "Fuerteventura", Country: "Spain", Latitude: 28.4527, Longitude: -13.8638, Aliases: []string{"fuerteventura", "puerto del rosario"}},
	"BIO": {IATA: "BIO", Name: "Bilbao Airport", City: "Bilbao", Country: "Spain", Latitude: 43.3011, Longitude: -2.9106, Aliases: []string{"bilbao", "basque"}},
	"SDR": {IATA: "SDR", Name: "Santander Airport", City: "Santander", Country: "Spain", Latitude: 43.4271, Longitude: -3.8200, Aliases: []string{"santander", "cantabria"}},
	"SCQ": {IATA: "SCQ", Name: "Santiago de Compostela Airport", City: "Santiago", Country: "Spain", Latitude: 42.8963, Longitude: -8.4151, Aliases: []string{"santiago", "galicia"}},

	// Italy
	"BGY": {IATA: "BGY", Name: "Milan Bergamo Airport (Orio al Serio)", City: "Milan", Country: "Italy", Latitude: 45.6739, Longitude: 9.7042, Aliases: []string{"milan", "bergamo", "lombardy"}},
	"MXP": {IATA: "MXP", Name: "Milan Malpensa Airport", City: "Milan", Country: "Italy", Latitude: 45.6301, Longitude: 8.7255, Aliases: []string{"milan", "malpensa"}},
	"CIA": {IATA: "CIA", Name: "Rome Ciampino Airport", City: "Rome", Country: "Italy", Latitude: 41.7994, Longitude: 12.5949, Aliases: []string{"rome", "ciampino", "lazio"}},
	"FCO": {IATA: "FCO", Name: "Rome Fiumicino Leonardo da Vinci Airport", City: "Rome", Country: "Italy", Latitude: 41.8003, Longitude: 12.2389, Aliases: []string{"rome", "fiumicino"}},
	"VCE": {IATA: "VCE", Name: "Venice Marco Polo Airport", City: "Venice", Country: "Italy", Latitude: 45.5053, Longitude: 12.3519, Aliases: []string{"venice", "venezia", "veneto"}},
	"TSF": {IATA: "TSF", Name: "Treviso Airport", City: "Treviso", Country: "Italy", Latitude: 45.6484, Longitude: 12.1944, Aliases: []string{"venice", "treviso"}},
	"BLQ": {IATA: "BLQ", Name: "Bologna Guglielmo Marconi Airport", City: "Bologna", Country: "Italy", Latitude: 44.5354, Longitude: 11.2887, Aliases: []string{"bologna", "emilia romagna"}},
	"PSA": {IATA: "PSA", Name: "Pisa International Airport", City: "Pisa", Country: "Italy", Latitude: 43.6839, Longitude: 10.3927, Aliases: []string{"pisa", "florence", "tuscany"}},
	"TRS": {IATA: "TRS", Name: "Trieste – Friuli Venezia Giulia Airport", City: "Trieste", Country: "Italy", Latitude: 45.8275, Longitude: 13.4722, Aliases: []string{"trieste", "friuli"}},
	"NAP": {IATA: "NAP", Name: "Naples International Airport", City: "Naples", Country: "Italy", Latitude: 40.8860, Longitude: 14.2908, Aliases: []string{"naples", "napoli", "amalfi", "campania"}},
	"BRI": {IATA: "BRI", Name: "Bari Karol Wojtyła Airport", City: "Bari", Country: "Italy", Latitude: 41.1389, Longitude: 16.7606, Aliases: []string{"bari", "puglia", "apulia", "polignano", "monopoli", "alberobello", "matera"}},
	"BDS": {IATA: "BDS", Name: "Brindisi Airport", City: "Brindisi", Country: "Italy", Latitude: 40.6576, Longitude: 17.9470, Aliases: []string{"brindisi", "salento", "puglia", "apulia", "lecce", "ostuni"}},
	"SUF": {IATA: "SUF", Name: "Lamezia Terme International Airport", City: "Lamezia Terme", Country: "Italy", Latitude: 38.9054, Longitude: 16.2423, Aliases: []string{"lamezia", "calabria"}},
	"CTA": {IATA: "CTA", Name: "Catania-Fontanarossa Airport", City: "Catania", Country: "Italy", Latitude: 37.4668, Longitude: 15.0664, Aliases: []string{"catania", "sicily", "sicilia", "etna", "taormina"}},
	"PMO": {IATA: "PMO", Name: "Palermo Airport", City: "Palermo", Country: "Italy", Latitude: 38.1760, Longitude: 13.0910, Aliases: []string{"palermo", "sicily", "sicilia", "mondello", "cefalù"}},
	"TPS": {IATA: "TPS", Name: "Trapani–Birgi Airport (Vincenzo Florio)", City: "Trapani", Country: "Italy", Latitude: 37.9069, Longitude: 12.4881, Aliases: []string{"trapani", "sicily", "sicilia", "marsala"}},
	"CAG": {IATA: "CAG", Name: "Cagliari Elmas Airport", City: "Cagliari", Country: "Italy", Latitude: 39.2515, Longitude: 9.0543, Aliases: []string{"cagliari", "sardinia"}},
	"TRN": {IATA: "TRN", Name: "Turin Airport", City: "Turin", Country: "Italy", Latitude: 45.2008, Longitude: 7.6496, Aliases: []string{"turin", "torino", "piedmont"}},
	"VRN": {IATA: "VRN", Name: "Verona Villafranca Airport", City: "Verona", Country: "Italy", Latitude: 45.3957, Longitude: 10.8885, Aliases: []string{"verona", "garda"}},

	// France
	"BVA": {IATA: "BVA", Name: "Beauvais–Tillé Airport", City: "Paris", Country: "France", Latitude: 49.4544, Longitude: 2.1128, Aliases: []string{"paris", "beauvais"}},
	"MRS": {IATA: "MRS", Name: "Marseille Provence Airport", City: "Marseille", Country: "France", Latitude: 43.4367, Longitude: 5.2150, Aliases: []string{"marseille", "provence"}},
	"NCE": {IATA: "NCE", Name: "Nice Côte d'Azur Airport", City: "Nice", Country: "France", Latitude: 43.6584, Longitude: 7.2159, Aliases: []string{"nice", "french riviera", "monaco"}},
	"TLS": {IATA: "TLS", Name: "Toulouse–Blagnac Airport", City: "Toulouse", Country: "France", Latitude: 43.6291, Longitude: 1.3638, Aliases: []string{"toulouse", "occitanie"}},
	"BOD": {IATA: "BOD", Name: "Bordeaux–Mérignac Airport", City: "Bordeaux", Country: "France", Latitude: 44.8283, Longitude: -0.7156, Aliases: []string{"bordeaux", "aquitaine"}},
	"NTE": {IATA: "NTE", Name: "Nantes Atlantique Airport", City: "Nantes", Country: "France", Latitude: 47.1532, Longitude: -1.6107, Aliases: []string{"nantes", "brittany"}},
	"LYS": {IATA: "LYS", Name: "Lyon–Saint-Exupéry Airport", City: "Lyon", Country: "France", Latitude: 45.7256, Longitude: 5.0811, Aliases: []string{"lyon", "rhone"}},

	// Germany, Austria, Switzerland
	"BER": {IATA: "BER", Name: "Berlin Brandenburg Airport", City: "Berlin", Country: "Germany", Latitude: 52.3667, Longitude: 13.5033, Aliases: []string{"berlin", "brandenburg"}},
	"FRA": {IATA: "FRA", Name: "Frankfurt Airport", City: "Frankfurt", Country: "Germany", Latitude: 50.0379, Longitude: 8.5622, Aliases: []string{"frankfurt", "main"}},
	"HHN": {IATA: "HHN", Name: "Frankfurt–Hahn Airport", City: "Frankfurt", Country: "Germany", Latitude: 49.9487, Longitude: 7.2639, Aliases: []string{"frankfurt hahn", "hahn"}},
	"CGN": {IATA: "CGN", Name: "Cologne Bonn Airport", City: "Cologne", Country: "Germany", Latitude: 50.8659, Longitude: 7.1427, Aliases: []string{"cologne", "koln", "bonn"}},
	"NRN": {IATA: "NRN", Name: "Weeze Airport (Düsseldorf Weeze)", City: "Weeze", Country: "Germany", Latitude: 51.6024, Longitude: 6.1422, Aliases: []string{"weeze", "dusseldorf weeze", "niederrhein", "germany", "netherlands"}},
	"FKB": {IATA: "FKB", Name: "Karlsruhe/Baden-Baden Airport", City: "Karlsruhe", Country: "Germany", Latitude: 48.7794, Longitude: 8.0805, Aliases: []string{"karlsruhe", "baden-baden", "black forest", "alsace"}},
	"FMM": {IATA: "FMM", Name: "Memmingen Airport", City: "Munich", Country: "Germany", Latitude: 47.9888, Longitude: 10.2395, Aliases: []string{"munich", "memmingen", "allgau", "bavaria"}},
	"HAM": {IATA: "HAM", Name: "Hamburg Airport", City: "Hamburg", Country: "Germany", Latitude: 53.6304, Longitude: 9.9882, Aliases: []string{"hamburg"}},
	"NUE": {IATA: "NUE", Name: "Nuremberg Airport", City: "Nuremberg", Country: "Germany", Latitude: 49.4987, Longitude: 11.0780, Aliases: []string{"nuremberg", "nurnberg", "bavaria"}},
	"DTM": {IATA: "DTM", Name: "Dortmund Airport", City: "Dortmund", Country: "Germany", Latitude: 51.5183, Longitude: 7.6122, Aliases: []string{"dortmund", "ruhr"}},
	"BRE": {IATA: "BRE", Name: "Bremen Airport", City: "Bremen", Country: "Germany", Latitude: 53.0475, Longitude: 8.7867, Aliases: []string{"bremen"}},
	"VIE": {IATA: "VIE", Name: "Vienna International Airport", City: "Vienna", Country: "Austria", Latitude: 48.1103, Longitude: 16.5697, Aliases: []string{"vienna", "wien", "austria"}},
	"SZG": {IATA: "SZG", Name: "Salzburg Airport", City: "Salzburg", Country: "Austria", Latitude: 47.7944, Longitude: 13.0043, Aliases: []string{"salzburg", "alps"}},
	"BSL": {IATA: "BSL", Name: "EuroAirport Basel Mulhouse Freiburg", City: "Basel", Country: "Switzerland", Latitude: 47.5900, Longitude: 7.5292, Aliases: []string{"basel", "mulhouse", "freiburg"}},

	// Central & Eastern Europe
	"BTS": {IATA: "BTS", Name: "Bratislava Airport", City: "Bratislava", Country: "Slovakia", Latitude: 48.1700, Longitude: 17.2127, Aliases: []string{"bratislava", "vienna east"}},
	"PRG": {IATA: "PRG", Name: "Václav Havel Airport Prague", City: "Prague", Country: "Czech Republic", Latitude: 50.1008, Longitude: 14.2600, Aliases: []string{"prague", "praha", "czechia"}},
	"BUD": {IATA: "BUD", Name: "Budapest Ferenc Liszt Airport", City: "Budapest", Country: "Hungary", Latitude: 47.4369, Longitude: 19.2556, Aliases: []string{"budapest", "hungary"}},
	"WAW": {IATA: "WAW", Name: "Warsaw Chopin Airport", City: "Warsaw", Country: "Poland", Latitude: 52.1672, Longitude: 20.9679, Aliases: []string{"warsaw", "warszawa"}},
	"WMI": {IATA: "WMI", Name: "Warsaw Modlin Airport", City: "Warsaw", Country: "Poland", Latitude: 52.4511, Longitude: 20.6517, Aliases: []string{"warsaw modlin", "modlin"}},
	"KRK": {IATA: "KRK", Name: "Kraków John Paul II Airport", City: "Kraków", Country: "Poland", Latitude: 50.0777, Longitude: 19.7848, Aliases: []string{"krakow", "cracow"}},
	"KTW": {IATA: "KTW", Name: "Katowice Airport", City: "Katowice", Country: "Poland", Latitude: 50.4743, Longitude: 19.0800, Aliases: []string{"katowice", "silesia"}},
	"GDN": {IATA: "GDN", Name: "Gdańsk Lech Wałęsa Airport", City: "Gdańsk", Country: "Poland", Latitude: 54.3776, Longitude: 18.4662, Aliases: []string{"gdansk", "baltic"}},
	"WRO": {IATA: "WRO", Name: "Wrocław Copernicus Airport", City: "Wrocław", Country: "Poland", Latitude: 51.1027, Longitude: 16.8858, Aliases: []string{"wroclaw", "breslau"}},

	// Scandinavia
	"CPH": {IATA: "CPH", Name: "Copenhagen Airport", City: "Copenhagen", Country: "Denmark", Latitude: 55.6180, Longitude: 12.6508, Aliases: []string{"copenhagen", "kobenhavn", "malmo"}},
	"BLL": {IATA: "BLL", Name: "Billund Airport", City: "Billund", Country: "Denmark", Latitude: 55.7403, Longitude: 9.1518, Aliases: []string{"billund", "legoland", "jutland"}},
	"AAL": {IATA: "AAL", Name: "Aalborg Airport", City: "Aalborg", Country: "Denmark", Latitude: 57.0928, Longitude: 9.8492, Aliases: []string{"aalborg"}},
	"ARN": {IATA: "ARN", Name: "Stockholm Arlanda Airport", City: "Stockholm", Country: "Sweden", Latitude: 59.6519, Longitude: 17.9186, Aliases: []string{"stockholm", "arlanda"}},
	"GOT": {IATA: "GOT", Name: "Göteborg Landvetter Airport", City: "Gothenburg", Country: "Sweden", Latitude: 57.6628, Longitude: 12.2798, Aliases: []string{"gothenburg", "goteborg"}},
	"OSL": {IATA: "OSL", Name: "Oslo Airport (Gardermoen)", City: "Oslo", Country: "Norway", Latitude: 60.1976, Longitude: 11.1004, Aliases: []string{"oslo", "gardermoen"}},
	"TRF": {IATA: "TRF", Name: "Sandefjord Airport (Torp)", City: "Oslo", Country: "Norway", Latitude: 59.1867, Longitude: 10.2586, Aliases: []string{"torp", "sandefjord", "oslo south"}},

	// Greece & Cyprus & Croatia & Malta
	"ATH": {IATA: "ATH", Name: "Athens International Airport", City: "Athens", Country: "Greece", Latitude: 37.9364, Longitude: 23.9445, Aliases: []string{"athens", "attica"}},
	"SKG": {IATA: "SKG", Name: "Thessaloniki Airport", City: "Thessaloniki", Country: "Greece", Latitude: 40.5197, Longitude: 22.9709, Aliases: []string{"thessaloniki", "macedonia"}},
	"CHQ": {IATA: "CHQ", Name: "Chania International Airport", City: "Chania", Country: "Greece", Latitude: 35.5317, Longitude: 24.1497, Aliases: []string{"chania", "crete"}},
	"HER": {IATA: "HER", Name: "Heraklion International Airport", City: "Heraklion", Country: "Greece", Latitude: 35.3397, Longitude: 25.1803, Aliases: []string{"heraklion", "crete"}},
	"RHO": {IATA: "RHO", Name: "Rhodes International Airport", City: "Rhodes", Country: "Greece", Latitude: 36.4054, Longitude: 28.0862, Aliases: []string{"rhodes", "dodecanese"}},
	"CFU": {IATA: "CFU", Name: "Corfu International Airport", City: "Corfu", Country: "Greece", Latitude: 39.6019, Longitude: 19.9117, Aliases: []string{"corfu", "kerkyra"}},
	"JTR": {IATA: "JTR", Name: "Santorini International Airport", City: "Santorini", Country: "Greece", Latitude: 36.3992, Longitude: 25.4793, Aliases: []string{"santorini", "thira"}},
	"JMK": {IATA: "JMK", Name: "Mykonos Airport", City: "Mykonos", Country: "Greece", Latitude: 37.4351, Longitude: 25.3481, Aliases: []string{"mykonos"}},
	"MLA": {IATA: "MLA", Name: "Malta International Airport", City: "Valletta", Country: "Malta", Latitude: 35.8575, Longitude: 14.4775, Aliases: []string{"malta", "luqa", "valletta"}},
	"PFO": {IATA: "PFO", Name: "Paphos International Airport", City: "Paphos", Country: "Cyprus", Latitude: 34.7180, Longitude: 32.4857, Aliases: []string{"paphos", "cyprus"}},
	"LCA": {IATA: "LCA", Name: "Larnaca International Airport", City: "Larnaca", Country: "Cyprus", Latitude: 34.8751, Longitude: 33.6249, Aliases: []string{"larnaca", "cyprus"}},
	"ZAG": {IATA: "ZAG", Name: "Zagreb Airport", City: "Zagreb", Country: "Croatia", Latitude: 45.7430, Longitude: 16.0688, Aliases: []string{"zagreb", "croatia"}},
	"ZAD": {IATA: "ZAD", Name: "Zadar Airport", City: "Zadar", Country: "Croatia", Latitude: 44.1083, Longitude: 15.3467, Aliases: []string{"zadar", "dalmatia"}},
	"SPU": {IATA: "SPU", Name: "Split Airport", City: "Split", Country: "Croatia", Latitude: 43.5389, Longitude: 16.2980, Aliases: []string{"split", "dalmatia"}},
	"DBV": {IATA: "DBV", Name: "Dubrovnik Airport", City: "Dubrovnik", Country: "Croatia", Latitude: 42.5614, Longitude: 18.2682, Aliases: []string{"dubrovnik"}},

	// Morocco
	"RAK": {IATA: "RAK", Name: "Marrakesh Menara Airport", City: "Marrakesh", Country: "Morocco", Latitude: 31.6069, Longitude: -8.0363, Aliases: []string{"marrakech", "marrakesh", "morocco"}},
	"AGA": {IATA: "AGA", Name: "Agadir–Al Massira Airport", City: "Agadir", Country: "Morocco", Latitude: 30.3250, Longitude: -9.4131, Aliases: []string{"agadir", "morocco"}},
	"FEZ": {IATA: "FEZ", Name: "Fès–Saïss Airport", City: "Fes", Country: "Morocco", Latitude: 33.9273, Longitude: -4.9780, Aliases: []string{"fez", "fes", "morocco"}},
	"TNG": {IATA: "TNG", Name: "Tangier Ibn Battouta Airport", City: "Tangier", Country: "Morocco", Latitude: 35.7269, Longitude: -5.9169, Aliases: []string{"tangier", "morocco"}},
}

// DistanceKm calculates the great-circle distance between two points in kilometers using the Haversine formula.
func DistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	rLat1 := lat1 * (math.Pi / 180.0)
	rLat2 := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return math.Round(earthRadiusKm*c*10) / 10
}

// FindAirport finds an airport by exact IATA code or fuzzy name/city/alias match.
func FindAirport(query string) (Airport, bool) {
	q := strings.TrimSpace(strings.ToUpper(query))
	if ap, exists := Directory[q]; exists {
		return ap, true
	}

	qLower := strings.ToLower(strings.TrimSpace(query))
	for _, ap := range Directory {
		if strings.ToLower(ap.City) == qLower || strings.ToLower(ap.Name) == qLower {
			return ap, true
		}
		for _, alias := range ap.Aliases {
			if alias == qLower {
				return ap, true
			}
		}
	}

	// Partial substring match
	for _, ap := range Directory {
		if strings.Contains(strings.ToLower(ap.City), qLower) ||
			strings.Contains(strings.ToLower(ap.Name), qLower) {
			return ap, true
		}
	}

	return Airport{}, false
}

// FindNearbyAirports returns airports located within maxRadiusKm from the given airport or city query.
func FindNearbyAirports(query string, maxRadiusKm float64, maxResults int) ([]NearbyAirport, Airport, bool) {
	anchor, found := FindAirport(query)
	if !found {
		return nil, Airport{}, false
	}

	if maxRadiusKm <= 0 {
		maxRadiusKm = 150.0 // Default 150 km radius
	}
	if maxResults <= 0 {
		maxResults = 10
	}

	var candidates []NearbyAirport
	for _, ap := range Directory {
		dist := DistanceKm(anchor.Latitude, anchor.Longitude, ap.Latitude, ap.Longitude)
		if dist <= maxRadiusKm {
			candidates = append(candidates, NearbyAirport{
				Airport:    ap,
				DistanceKm: dist,
			})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].DistanceKm < candidates[j].DistanceKm
	})

	if len(candidates) > maxResults {
		candidates = candidates[:maxResults]
	}

	return candidates, anchor, true
}

// ResolveAirportCodes takes an input identifier (IATA, city, region) and returns a list of unique IATA codes,
// starting with the primary airport and followed by nearby airports within maxRadiusKm.
func ResolveAirportCodes(query string, maxRadiusKm float64) []string {
	nearby, anchor, found := FindNearbyAirports(query, maxRadiusKm, 6)
	if !found {
		// Fallback: if query is a 3-letter IATA code, return it as-is
		clean := strings.ToUpper(strings.TrimSpace(query))
		if len(clean) == 3 {
			return []string{clean}
		}
		return nil
	}

	seen := make(map[string]bool)
	var result []string

	// Anchor first
	result = append(result, anchor.IATA)
	seen[anchor.IATA] = true

	for _, n := range nearby {
		if !seen[n.IATA] {
			result = append(result, n.IATA)
			seen[n.IATA] = true
		}
	}

	return result
}
