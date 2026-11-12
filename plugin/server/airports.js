'use strict';

// Curated directory of European and North African airports served by Ryanair.
// Ported from internal/airports/airports.go. Coordinates are decimal degrees.

const DIRECTORY = {
  // Ireland
  DUB: { iata: 'DUB', name: 'Dublin Airport', city: 'Dublin', country: 'Ireland', latitude: 53.4213, longitude: -6.2701, aliases: ['dublin', 'leinster'] },
  SNN: { iata: 'SNN', name: 'Shannon Airport', city: 'Shannon', country: 'Ireland', latitude: 52.702, longitude: -8.9248, aliases: ['limerick', 'galway', 'shannon'] },
  ORK: { iata: 'ORK', name: 'Cork Airport', city: 'Cork', country: 'Ireland', latitude: 51.8413, longitude: -8.4911, aliases: ['cork', 'munster'] },
  NOC: { iata: 'NOC', name: 'Ireland West Airport Knock', city: 'Knock', country: 'Ireland', latitude: 53.9103, longitude: -8.8185, aliases: ['knock', 'mayo', 'connacht'] },
  BFS: { iata: 'BFS', name: 'Belfast International Airport', city: 'Belfast', country: 'United Kingdom', latitude: 54.6575, longitude: -6.2158, aliases: ['belfast', 'northern ireland'] },

  // United Kingdom
  STN: { iata: 'STN', name: 'London Stansted Airport', city: 'London', country: 'United Kingdom', latitude: 51.886, longitude: 0.2389, aliases: ['london', 'stansted', 'essex'] },
  LGW: { iata: 'LGW', name: 'London Gatwick Airport', city: 'London', country: 'United Kingdom', latitude: 51.1537, longitude: -0.1821, aliases: ['london', 'gatwick', 'surrey'] },
  LTN: { iata: 'LTN', name: 'London Luton Airport', city: 'London', country: 'United Kingdom', latitude: 51.8747, longitude: -0.3683, aliases: ['london', 'luton', 'bedfordshire'] },
  MAN: { iata: 'MAN', name: 'Manchester Airport', city: 'Manchester', country: 'United Kingdom', latitude: 53.3537, longitude: -2.275, aliases: ['manchester', 'north west'] },
  BHX: { iata: 'BHX', name: 'Birmingham Airport', city: 'Birmingham', country: 'United Kingdom', latitude: 52.4524, longitude: -1.7435, aliases: ['birmingham', 'midlands'] },
  EDI: { iata: 'EDI', name: 'Edinburgh Airport', city: 'Edinburgh', country: 'United Kingdom', latitude: 55.9508, longitude: -3.3725, aliases: ['edinburgh', 'scotland'] },
  GLA: { iata: 'GLA', name: 'Glasgow Airport', city: 'Glasgow', country: 'United Kingdom', latitude: 55.8719, longitude: -4.4331, aliases: ['glasgow', 'scotland'] },
  BRS: { iata: 'BRS', name: 'Bristol Airport', city: 'Bristol', country: 'United Kingdom', latitude: 51.3827, longitude: -2.7191, aliases: ['bristol', 'bath', 'somerset'] },
  LPL: { iata: 'LPL', name: 'Liverpool John Lennon Airport', city: 'Liverpool', country: 'United Kingdom', latitude: 53.3336, longitude: -2.8497, aliases: ['liverpool', 'merseyside'] },
  EMA: { iata: 'EMA', name: 'East Midlands Airport', city: 'Nottingham', country: 'United Kingdom', latitude: 52.8311, longitude: -1.3281, aliases: ['nottingham', 'derby', 'leicester'] },
  CWL: { iata: 'CWL', name: 'Cardiff Airport', city: 'Cardiff', country: 'United Kingdom', latitude: 51.3967, longitude: -3.3433, aliases: ['cardiff', 'wales'] },
  NCL: { iata: 'NCL', name: 'Newcastle Airport', city: 'Newcastle', country: 'United Kingdom', latitude: 55.0375, longitude: -1.6917, aliases: ['newcastle', 'tyne'] },
  LBA: { iata: 'LBA', name: 'Leeds Bradford Airport', city: 'Leeds', country: 'United Kingdom', latitude: 53.8659, longitude: -1.6606, aliases: ['leeds', 'bradford', 'yorkshire'] },

  // Belgium, Luxembourg, Netherlands, northern France
  CRL: { iata: 'CRL', name: 'Brussels South Charleroi Airport', city: 'Brussels', country: 'Belgium', latitude: 50.4592, longitude: 4.4538, aliases: ['brussels', 'charleroi', 'wallonia', 'belgium'] },
  BRU: { iata: 'BRU', name: 'Brussels Airport (Zaventem)', city: 'Brussels', country: 'Belgium', latitude: 50.9014, longitude: 4.4844, aliases: ['brussels', 'zaventem', 'belgium'] },
  LGG: { iata: 'LGG', name: 'Liege Airport', city: 'Liege', country: 'Belgium', latitude: 50.6374, longitude: 5.4432, aliases: ['liege', 'luik', 'belgium'] },
  LUX: { iata: 'LUX', name: 'Luxembourg Airport', city: 'Luxembourg', country: 'Luxembourg', latitude: 49.6233, longitude: 6.2044, aliases: ['luxembourg', 'lux'] },
  EIN: { iata: 'EIN', name: 'Eindhoven Airport', city: 'Eindhoven', country: 'Netherlands', latitude: 51.4501, longitude: 5.3745, aliases: ['eindhoven', 'brabant', 'netherlands'] },
  AMS: { iata: 'AMS', name: 'Amsterdam Airport Schiphol', city: 'Amsterdam', country: 'Netherlands', latitude: 52.3105, longitude: 4.7683, aliases: ['amsterdam', 'schiphol', 'holland'] },
  MST: { iata: 'MST', name: 'Maastricht Aachen Airport', city: 'Maastricht', country: 'Netherlands', latitude: 50.9117, longitude: 5.7701, aliases: ['maastricht', 'aachen', 'limburg'] },
  LIL: { iata: 'LIL', name: 'Lille Airport', city: 'Lille', country: 'France', latitude: 50.5619, longitude: 3.0894, aliases: ['lille', 'flanders'] },

  // Portugal
  FAO: { iata: 'FAO', name: 'Faro Airport', city: 'Faro', country: 'Portugal', latitude: 37.0176, longitude: -7.9659, aliases: ['faro', 'algarve', 'portugal', 'albufeira', 'lagos', 'vilamoura'] },
  LIS: { iata: 'LIS', name: 'Humberto Delgado Airport (Lisbon)', city: 'Lisbon', country: 'Portugal', latitude: 38.7742, longitude: -9.1342, aliases: ['lisbon', 'lisboa', 'sintra', 'cascais'] },
  OPO: { iata: 'OPO', name: 'Francisco Sa Carneiro Airport (Porto)', city: 'Porto', country: 'Portugal', latitude: 41.2481, longitude: -8.6814, aliases: ['porto', 'oporto', 'douro'] },
  FNC: { iata: 'FNC', name: 'Cristiano Ronaldo International Airport (Madeira)', city: 'Funchal', country: 'Portugal', latitude: 32.6979, longitude: -16.7745, aliases: ['madeira', 'funchal'] },
  PDL: { iata: 'PDL', name: 'Joao Paulo II Airport (Azores)', city: 'Ponta Delgada', country: 'Portugal', latitude: 37.7412, longitude: -25.6979, aliases: ['azores', 'ponta delgada', 'sao miguel'] },

  // Spain
  MAD: { iata: 'MAD', name: 'Adolfo Suarez Madrid-Barajas Airport', city: 'Madrid', country: 'Spain', latitude: 40.4839, longitude: -3.568, aliases: ['madrid', 'barajas', 'spain'] },
  BCN: { iata: 'BCN', name: 'Josep Tarradellas Barcelona-El Prat Airport', city: 'Barcelona', country: 'Spain', latitude: 41.2974, longitude: 2.0833, aliases: ['barcelona', 'el prat', 'catalonia'] },
  GRO: { iata: 'GRO', name: 'Girona-Costa Brava Airport', city: 'Girona', country: 'Spain', latitude: 41.901, longitude: 2.7606, aliases: ['girona', 'costa brava', 'barcelona north'] },
  REU: { iata: 'REU', name: 'Reus Airport', city: 'Reus', country: 'Spain', latitude: 41.1474, longitude: 1.1672, aliases: ['reus', 'tarragona', 'salou', 'costa daurada'] },
  AGP: { iata: 'AGP', name: 'Malaga-Costa del Sol Airport', city: 'Malaga', country: 'Spain', latitude: 36.6749, longitude: -4.4991, aliases: ['malaga', 'costa del sol', 'marbella', 'torremolinos', 'andalucia', 'andalusia', 'tarifa', 'gibraltar'] },
  SVQ: { iata: 'SVQ', name: 'Seville Airport', city: 'Seville', country: 'Spain', latitude: 37.418, longitude: -5.8931, aliases: ['seville', 'sevilla', 'andalucia', 'andalusia'] },
  XRY: { iata: 'XRY', name: 'Jerez Airport', city: 'Jerez de la Frontera', country: 'Spain', latitude: 36.7447, longitude: -6.0601, aliases: ['jerez', 'cadiz', 'tarifa', 'gibraltar', 'andalucia', 'andalusia'] },
  ALC: { iata: 'ALC', name: 'Alicante-Elche Miguel Hernandez Airport', city: 'Alicante', country: 'Spain', latitude: 38.2822, longitude: -0.5582, aliases: ['alicante', 'costa blanca', 'benidorm', 'altea', 'calpe', 'elche'] },
  RMU: { iata: 'RMU', name: 'Region de Murcia International Airport', city: 'Murcia', country: 'Spain', latitude: 37.803, longitude: -1.125, aliases: ['murcia', 'corvera', 'costa calida'] },
  VLC: { iata: 'VLC', name: 'Valencia Airport', city: 'Valencia', country: 'Spain', latitude: 39.4893, longitude: -0.4816, aliases: ['valencia'] },
  ZAZ: { iata: 'ZAZ', name: 'Zaragoza Airport', city: 'Zaragoza', country: 'Spain', latitude: 41.6662, longitude: -1.0416, aliases: ['zaragoza', 'aragon'] },
  PMI: { iata: 'PMI', name: 'Palma de Mallorca Airport', city: 'Palma de Mallorca', country: 'Spain', latitude: 39.5517, longitude: 2.7388, aliases: ['mallorca', 'majorca', 'palma', 'balearics'] },
  IBZ: { iata: 'IBZ', name: 'Ibiza Airport', city: 'Ibiza', country: 'Spain', latitude: 38.8729, longitude: 1.3731, aliases: ['ibiza', 'eivissa', 'balearics'] },
  MAH: { iata: 'MAH', name: 'Menorca Airport', city: 'Menorca', country: 'Spain', latitude: 39.8626, longitude: 4.2186, aliases: ['menorca', 'minorca', 'mahon'] },
  TFS: { iata: 'TFS', name: 'Tenerife South Airport', city: 'Tenerife', country: 'Spain', latitude: 28.0445, longitude: -16.5725, aliases: ['tenerife', 'canaries', 'reina sofia'] },
  TFN: { iata: 'TFN', name: 'Tenerife North Airport', city: 'Tenerife', country: 'Spain', latitude: 28.4827, longitude: -16.3415, aliases: ['tenerife north', 'los rodeos'] },
  LPA: { iata: 'LPA', name: 'Gran Canaria Airport', city: 'Gran Canaria', country: 'Spain', latitude: 27.9319, longitude: -15.3866, aliases: ['gran canaria', 'las palmas', 'canaries'] },
  ACE: { iata: 'ACE', name: 'Lanzarote Airport', city: 'Lanzarote', country: 'Spain', latitude: 28.9455, longitude: -13.6052, aliases: ['lanzarote', 'arrecife', 'canaries'] },
  FUE: { iata: 'FUE', name: 'Fuerteventura Airport', city: 'Fuerteventura', country: 'Spain', latitude: 28.4527, longitude: -13.8638, aliases: ['fuerteventura', 'puerto del rosario'] },
  BIO: { iata: 'BIO', name: 'Bilbao Airport', city: 'Bilbao', country: 'Spain', latitude: 43.3011, longitude: -2.9106, aliases: ['bilbao', 'basque'] },
  SDR: { iata: 'SDR', name: 'Santander Airport', city: 'Santander', country: 'Spain', latitude: 43.4271, longitude: -3.82, aliases: ['santander', 'cantabria'] },
  SCQ: { iata: 'SCQ', name: 'Santiago de Compostela Airport', city: 'Santiago', country: 'Spain', latitude: 42.8963, longitude: -8.4151, aliases: ['santiago', 'galicia'] },

  // Italy
  BGY: { iata: 'BGY', name: 'Milan Bergamo Airport (Orio al Serio)', city: 'Milan', country: 'Italy', latitude: 45.6739, longitude: 9.7042, aliases: ['milan', 'bergamo', 'lombardy'] },
  MXP: { iata: 'MXP', name: 'Milan Malpensa Airport', city: 'Milan', country: 'Italy', latitude: 45.6301, longitude: 8.7255, aliases: ['milan', 'malpensa'] },
  CIA: { iata: 'CIA', name: 'Rome Ciampino Airport', city: 'Rome', country: 'Italy', latitude: 41.7994, longitude: 12.5949, aliases: ['rome', 'ciampino', 'lazio'] },
  FCO: { iata: 'FCO', name: 'Rome Fiumicino Leonardo da Vinci Airport', city: 'Rome', country: 'Italy', latitude: 41.8003, longitude: 12.2389, aliases: ['rome', 'fiumicino'] },
  VCE: { iata: 'VCE', name: 'Venice Marco Polo Airport', city: 'Venice', country: 'Italy', latitude: 45.5053, longitude: 12.3519, aliases: ['venice', 'venezia', 'veneto'] },
  TSF: { iata: 'TSF', name: 'Treviso Airport', city: 'Treviso', country: 'Italy', latitude: 45.6484, longitude: 12.1944, aliases: ['venice', 'treviso'] },
  BLQ: { iata: 'BLQ', name: 'Bologna Guglielmo Marconi Airport', city: 'Bologna', country: 'Italy', latitude: 44.5354, longitude: 11.2887, aliases: ['bologna', 'emilia romagna'] },
  PSA: { iata: 'PSA', name: 'Pisa International Airport', city: 'Pisa', country: 'Italy', latitude: 43.6839, longitude: 10.3927, aliases: ['pisa', 'florence', 'tuscany'] },
  TRS: { iata: 'TRS', name: 'Trieste - Friuli Venezia Giulia Airport', city: 'Trieste', country: 'Italy', latitude: 45.8275, longitude: 13.4722, aliases: ['trieste', 'friuli'] },
  NAP: { iata: 'NAP', name: 'Naples International Airport', city: 'Naples', country: 'Italy', latitude: 40.886, longitude: 14.2908, aliases: ['naples', 'napoli', 'amalfi', 'campania'] },
  BRI: { iata: 'BRI', name: 'Bari Karol Wojtyla Airport', city: 'Bari', country: 'Italy', latitude: 41.1389, longitude: 16.7606, aliases: ['bari', 'puglia', 'apulia', 'polignano', 'monopoli', 'alberobello', 'matera'] },
  BDS: { iata: 'BDS', name: 'Brindisi Airport', city: 'Brindisi', country: 'Italy', latitude: 40.6576, longitude: 17.947, aliases: ['brindisi', 'salento', 'puglia', 'apulia', 'lecce', 'ostuni'] },
  SUF: { iata: 'SUF', name: 'Lamezia Terme International Airport', city: 'Lamezia Terme', country: 'Italy', latitude: 38.9054, longitude: 16.2423, aliases: ['lamezia', 'calabria'] },
  CTA: { iata: 'CTA', name: 'Catania-Fontanarossa Airport', city: 'Catania', country: 'Italy', latitude: 37.4668, longitude: 15.0664, aliases: ['catania', 'sicily', 'sicilia', 'etna', 'taormina'] },
  PMO: { iata: 'PMO', name: 'Palermo Airport', city: 'Palermo', country: 'Italy', latitude: 38.176, longitude: 13.091, aliases: ['palermo', 'sicily', 'sicilia', 'mondello', 'cefalu'] },
  TPS: { iata: 'TPS', name: 'Trapani-Birgi Airport (Vincenzo Florio)', city: 'Trapani', country: 'Italy', latitude: 37.9069, longitude: 12.4881, aliases: ['trapani', 'sicily', 'sicilia', 'marsala'] },
  CAG: { iata: 'CAG', name: 'Cagliari Elmas Airport', city: 'Cagliari', country: 'Italy', latitude: 39.2515, longitude: 9.0543, aliases: ['cagliari', 'sardinia'] },
  TRN: { iata: 'TRN', name: 'Turin Airport', city: 'Turin', country: 'Italy', latitude: 45.2008, longitude: 7.6496, aliases: ['turin', 'torino', 'piedmont'] },
  VRN: { iata: 'VRN', name: 'Verona Villafranca Airport', city: 'Verona', country: 'Italy', latitude: 45.3957, longitude: 10.8885, aliases: ['verona', 'garda'] },

  // France
  BVA: { iata: 'BVA', name: 'Beauvais-Tille Airport', city: 'Paris', country: 'France', latitude: 49.4544, longitude: 2.1128, aliases: ['paris', 'beauvais'] },
  MRS: { iata: 'MRS', name: 'Marseille Provence Airport', city: 'Marseille', country: 'France', latitude: 43.4367, longitude: 5.215, aliases: ['marseille', 'provence'] },
  NCE: { iata: 'NCE', name: "Nice Cote d'Azur Airport", city: 'Nice', country: 'France', latitude: 43.6584, longitude: 7.2159, aliases: ['nice', 'french riviera', 'monaco'] },
  TLS: { iata: 'TLS', name: 'Toulouse-Blagnac Airport', city: 'Toulouse', country: 'France', latitude: 43.6291, longitude: 1.3638, aliases: ['toulouse', 'occitanie'] },
  BOD: { iata: 'BOD', name: 'Bordeaux-Merignac Airport', city: 'Bordeaux', country: 'France', latitude: 44.8283, longitude: -0.7156, aliases: ['bordeaux', 'aquitaine'] },
  NTE: { iata: 'NTE', name: 'Nantes Atlantique Airport', city: 'Nantes', country: 'France', latitude: 47.1532, longitude: -1.6107, aliases: ['nantes', 'brittany'] },
  LYS: { iata: 'LYS', name: 'Lyon-Saint-Exupery Airport', city: 'Lyon', country: 'France', latitude: 45.7256, longitude: 5.0811, aliases: ['lyon', 'rhone'] },

  // Germany, Austria, Switzerland
  BER: { iata: 'BER', name: 'Berlin Brandenburg Airport', city: 'Berlin', country: 'Germany', latitude: 52.3667, longitude: 13.5033, aliases: ['berlin', 'brandenburg'] },
  FRA: { iata: 'FRA', name: 'Frankfurt Airport', city: 'Frankfurt', country: 'Germany', latitude: 50.0379, longitude: 8.5622, aliases: ['frankfurt', 'main'] },
  HHN: { iata: 'HHN', name: 'Frankfurt-Hahn Airport', city: 'Frankfurt', country: 'Germany', latitude: 49.9487, longitude: 7.2639, aliases: ['frankfurt hahn', 'hahn'] },
  CGN: { iata: 'CGN', name: 'Cologne Bonn Airport', city: 'Cologne', country: 'Germany', latitude: 50.8659, longitude: 7.1427, aliases: ['cologne', 'koln', 'bonn'] },
  NRN: { iata: 'NRN', name: 'Weeze Airport (Dusseldorf Weeze)', city: 'Weeze', country: 'Germany', latitude: 51.6024, longitude: 6.1422, aliases: ['weeze', 'dusseldorf weeze', 'niederrhein', 'germany', 'netherlands'] },
  FKB: { iata: 'FKB', name: 'Karlsruhe/Baden-Baden Airport', city: 'Karlsruhe', country: 'Germany', latitude: 48.7794, longitude: 8.0805, aliases: ['karlsruhe', 'baden-baden', 'black forest', 'alsace'] },
  FMM: { iata: 'FMM', name: 'Memmingen Airport', city: 'Munich', country: 'Germany', latitude: 47.9888, longitude: 10.2395, aliases: ['munich', 'memmingen', 'allgau', 'bavaria'] },
  HAM: { iata: 'HAM', name: 'Hamburg Airport', city: 'Hamburg', country: 'Germany', latitude: 53.6304, longitude: 9.9882, aliases: ['hamburg'] },
  NUE: { iata: 'NUE', name: 'Nuremberg Airport', city: 'Nuremberg', country: 'Germany', latitude: 49.4987, longitude: 11.078, aliases: ['nuremberg', 'nurnberg', 'bavaria'] },
  DTM: { iata: 'DTM', name: 'Dortmund Airport', city: 'Dortmund', country: 'Germany', latitude: 51.5183, longitude: 7.6122, aliases: ['dortmund', 'ruhr'] },
  BRE: { iata: 'BRE', name: 'Bremen Airport', city: 'Bremen', country: 'Germany', latitude: 53.0475, longitude: 8.7867, aliases: ['bremen'] },
  VIE: { iata: 'VIE', name: 'Vienna International Airport', city: 'Vienna', country: 'Austria', latitude: 48.1103, longitude: 16.5697, aliases: ['vienna', 'wien', 'austria'] },
  SZG: { iata: 'SZG', name: 'Salzburg Airport', city: 'Salzburg', country: 'Austria', latitude: 47.7944, longitude: 13.0043, aliases: ['salzburg', 'alps'] },
  BSL: { iata: 'BSL', name: 'EuroAirport Basel Mulhouse Freiburg', city: 'Basel', country: 'Switzerland', latitude: 47.59, longitude: 7.5292, aliases: ['basel', 'mulhouse', 'freiburg'] },

  // Central and Eastern Europe
  BTS: { iata: 'BTS', name: 'Bratislava Airport', city: 'Bratislava', country: 'Slovakia', latitude: 48.17, longitude: 17.2127, aliases: ['bratislava', 'vienna east'] },
  PRG: { iata: 'PRG', name: 'Vaclav Havel Airport Prague', city: 'Prague', country: 'Czech Republic', latitude: 50.1008, longitude: 14.26, aliases: ['prague', 'praha', 'czechia'] },
  BUD: { iata: 'BUD', name: 'Budapest Ferenc Liszt Airport', city: 'Budapest', country: 'Hungary', latitude: 47.4369, longitude: 19.2556, aliases: ['budapest', 'hungary'] },
  WAW: { iata: 'WAW', name: 'Warsaw Chopin Airport', city: 'Warsaw', country: 'Poland', latitude: 52.1672, longitude: 20.9679, aliases: ['warsaw', 'warszawa'] },
  WMI: { iata: 'WMI', name: 'Warsaw Modlin Airport', city: 'Warsaw', country: 'Poland', latitude: 52.4511, longitude: 20.6517, aliases: ['warsaw modlin', 'modlin'] },
  KRK: { iata: 'KRK', name: 'Krakow John Paul II Airport', city: 'Krakow', country: 'Poland', latitude: 50.0777, longitude: 19.7848, aliases: ['krakow', 'cracow'] },
  KTW: { iata: 'KTW', name: 'Katowice Airport', city: 'Katowice', country: 'Poland', latitude: 50.4743, longitude: 19.08, aliases: ['katowice', 'silesia'] },
  GDN: { iata: 'GDN', name: 'Gdansk Lech Walesa Airport', city: 'Gdansk', country: 'Poland', latitude: 54.3776, longitude: 18.4662, aliases: ['gdansk', 'baltic'] },
  WRO: { iata: 'WRO', name: 'Wroclaw Copernicus Airport', city: 'Wroclaw', country: 'Poland', latitude: 51.1027, longitude: 16.8858, aliases: ['wroclaw', 'breslau'] },

  // Scandinavia
  CPH: { iata: 'CPH', name: 'Copenhagen Airport', city: 'Copenhagen', country: 'Denmark', latitude: 55.618, longitude: 12.6508, aliases: ['copenhagen', 'kobenhavn', 'malmo'] },
  BLL: { iata: 'BLL', name: 'Billund Airport', city: 'Billund', country: 'Denmark', latitude: 55.7403, longitude: 9.1518, aliases: ['billund', 'legoland', 'jutland'] },
  AAL: { iata: 'AAL', name: 'Aalborg Airport', city: 'Aalborg', country: 'Denmark', latitude: 57.0928, longitude: 9.8492, aliases: ['aalborg'] },
  ARN: { iata: 'ARN', name: 'Stockholm Arlanda Airport', city: 'Stockholm', country: 'Sweden', latitude: 59.6519, longitude: 17.9186, aliases: ['stockholm', 'arlanda'] },
  GOT: { iata: 'GOT', name: 'Goteborg Landvetter Airport', city: 'Gothenburg', country: 'Sweden', latitude: 57.6628, longitude: 12.2798, aliases: ['gothenburg', 'goteborg'] },
  OSL: { iata: 'OSL', name: 'Oslo Airport (Gardermoen)', city: 'Oslo', country: 'Norway', latitude: 60.1976, longitude: 11.1004, aliases: ['oslo', 'gardermoen'] },
  TRF: { iata: 'TRF', name: 'Sandefjord Airport (Torp)', city: 'Oslo', country: 'Norway', latitude: 59.1867, longitude: 10.2586, aliases: ['torp', 'sandefjord', 'oslo south'] },

  // Greece, Cyprus, Croatia, Malta
  ATH: { iata: 'ATH', name: 'Athens International Airport', city: 'Athens', country: 'Greece', latitude: 37.9364, longitude: 23.9445, aliases: ['athens', 'attica'] },
  SKG: { iata: 'SKG', name: 'Thessaloniki Airport', city: 'Thessaloniki', country: 'Greece', latitude: 40.5197, longitude: 22.9709, aliases: ['thessaloniki', 'macedonia'] },
  CHQ: { iata: 'CHQ', name: 'Chania International Airport', city: 'Chania', country: 'Greece', latitude: 35.5317, longitude: 24.1497, aliases: ['chania', 'crete'] },
  HER: { iata: 'HER', name: 'Heraklion International Airport', city: 'Heraklion', country: 'Greece', latitude: 35.3397, longitude: 25.1803, aliases: ['heraklion', 'crete'] },
  RHO: { iata: 'RHO', name: 'Rhodes International Airport', city: 'Rhodes', country: 'Greece', latitude: 36.4054, longitude: 28.0862, aliases: ['rhodes', 'dodecanese'] },
  CFU: { iata: 'CFU', name: 'Corfu International Airport', city: 'Corfu', country: 'Greece', latitude: 39.6019, longitude: 19.9117, aliases: ['corfu', 'kerkyra'] },
  JTR: { iata: 'JTR', name: 'Santorini International Airport', city: 'Santorini', country: 'Greece', latitude: 36.3992, longitude: 25.4793, aliases: ['santorini', 'thira'] },
  JMK: { iata: 'JMK', name: 'Mykonos Airport', city: 'Mykonos', country: 'Greece', latitude: 37.4351, longitude: 25.3481, aliases: ['mykonos'] },
  MLA: { iata: 'MLA', name: 'Malta International Airport', city: 'Valletta', country: 'Malta', latitude: 35.8575, longitude: 14.4775, aliases: ['malta', 'luqa', 'valletta'] },
  PFO: { iata: 'PFO', name: 'Paphos International Airport', city: 'Paphos', country: 'Cyprus', latitude: 34.718, longitude: 32.4857, aliases: ['paphos', 'cyprus'] },
  LCA: { iata: 'LCA', name: 'Larnaca International Airport', city: 'Larnaca', country: 'Cyprus', latitude: 34.8751, longitude: 33.6249, aliases: ['larnaca', 'cyprus'] },
  ZAG: { iata: 'ZAG', name: 'Zagreb Airport', city: 'Zagreb', country: 'Croatia', latitude: 45.743, longitude: 16.0688, aliases: ['zagreb', 'croatia'] },
  ZAD: { iata: 'ZAD', name: 'Zadar Airport', city: 'Zadar', country: 'Croatia', latitude: 44.1083, longitude: 15.3467, aliases: ['zadar', 'dalmatia'] },
  SPU: { iata: 'SPU', name: 'Split Airport', city: 'Split', country: 'Croatia', latitude: 43.5389, longitude: 16.298, aliases: ['split', 'dalmatia'] },
  DBV: { iata: 'DBV', name: 'Dubrovnik Airport', city: 'Dubrovnik', country: 'Croatia', latitude: 42.5614, longitude: 18.2682, aliases: ['dubrovnik'] },

  // Morocco
  RAK: { iata: 'RAK', name: 'Marrakesh Menara Airport', city: 'Marrakesh', country: 'Morocco', latitude: 31.6069, longitude: -8.0363, aliases: ['marrakech', 'marrakesh', 'morocco'] },
  AGA: { iata: 'AGA', name: 'Agadir-Al Massira Airport', city: 'Agadir', country: 'Morocco', latitude: 30.325, longitude: -9.4131, aliases: ['agadir', 'morocco'] },
  FEZ: { iata: 'FEZ', name: 'Fes-Saiss Airport', city: 'Fes', country: 'Morocco', latitude: 33.9273, longitude: -4.978, aliases: ['fez', 'fes', 'morocco'] },
  TNG: { iata: 'TNG', name: 'Tangier Ibn Battouta Airport', city: 'Tangier', country: 'Morocco', latitude: 35.7269, longitude: -5.9169, aliases: ['tangier', 'morocco'] },
};

const EARTH_RADIUS_KM = 6371.0;
const DEG_TO_RAD = Math.PI / 180;

/**
 * Great-circle distance between two coordinates, in kilometres,
 * rounded to one decimal place (Haversine formula).
 */
function distanceKm(lat1, lon1, lat2, lon2) {
  const dLat = (lat2 - lat1) * DEG_TO_RAD;
  const dLon = (lon2 - lon1) * DEG_TO_RAD;
  const rLat1 = lat1 * DEG_TO_RAD;
  const rLat2 = lat2 * DEG_TO_RAD;

  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(rLat1) * Math.cos(rLat2) * Math.sin(dLon / 2) * Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

  return Math.round(EARTH_RADIUS_KM * c * 10) / 10;
}

/**
 * Resolve a query to a single airport by exact IATA code, then by exact
 * city/name/alias match, then by partial city/name substring match.
 * Returns null when nothing matches.
 */
function findAirport(query) {
  if (typeof query !== 'string' || query.trim() === '') return null;

  const upper = query.trim().toUpperCase();
  if (Object.prototype.hasOwnProperty.call(DIRECTORY, upper)) {
    return DIRECTORY[upper];
  }

  const lower = query.trim().toLowerCase();
  const all = Object.values(DIRECTORY);

  for (const ap of all) {
    if (ap.city.toLowerCase() === lower || ap.name.toLowerCase() === lower) return ap;
    if (ap.aliases && ap.aliases.includes(lower)) return ap;
  }

  for (const ap of all) {
    if (ap.city.toLowerCase().includes(lower) || ap.name.toLowerCase().includes(lower)) return ap;
  }

  return null;
}

/**
 * Find airports within maxRadiusKm of the airport that `query` resolves to.
 * Returns { anchor, nearby } or null when the query resolves to no airport.
 * The anchor itself is included in `nearby` at distance 0.
 */
function findNearbyAirports(query, maxRadiusKm, maxResults) {
  const anchor = findAirport(query);
  if (!anchor) return null;

  const radius = maxRadiusKm > 0 ? maxRadiusKm : 150;
  const limit = maxResults > 0 ? maxResults : 10;

  const candidates = [];
  for (const ap of Object.values(DIRECTORY)) {
    const dist = distanceKm(anchor.latitude, anchor.longitude, ap.latitude, ap.longitude);
    if (dist <= radius) {
      candidates.push({
        iata: ap.iata,
        name: ap.name,
        city: ap.city,
        country: ap.country,
        latitude: ap.latitude,
        longitude: ap.longitude,
        distance_km: dist,
      });
    }
  }

  candidates.sort((a, b) => a.distance_km - b.distance_km);

  return { anchor, nearby: candidates.slice(0, limit) };
}

module.exports = { DIRECTORY, distanceKm, findAirport, findNearbyAirports };
