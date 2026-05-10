package validation

// AWGData holds electrical properties for a given AWG gauge.
type AWGData struct {
	ResistancePerMeter float64 // ohms per meter (copper, ~20°C)
	Ampacity           float64 // chassis wiring ampacity (typical marine/automotive)
}

// AWGTable maps AWG gauge numbers to electrical data.
// Resistance values: copper at 20°C. Ampacity: typical marine/chassis wiring ratings.
// Gauge "0" maps to 0, "00" to -1, "000" to -2, "0000" to -3.
var AWGTable = map[int]AWGData{
	16: {ResistancePerMeter: 0.01320, Ampacity: 25},
	14: {ResistancePerMeter: 0.00829, Ampacity: 35},
	12: {ResistancePerMeter: 0.00521, Ampacity: 45},
	10: {ResistancePerMeter: 0.00328, Ampacity: 60},
	8:  {ResistancePerMeter: 0.00206, Ampacity: 80},
	6:  {ResistancePerMeter: 0.00130, Ampacity: 120},
	4:  {ResistancePerMeter: 0.000815, Ampacity: 160},
	2:  {ResistancePerMeter: 0.000512, Ampacity: 210},
	1:  {ResistancePerMeter: 0.000406, Ampacity: 245},
	0:  {ResistancePerMeter: 0.000323, Ampacity: 285},
	-1: {ResistancePerMeter: 0.000256, Ampacity: 330},  // 2/0
	-2: {ResistancePerMeter: 0.000203, Ampacity: 385},  // 3/0
	-3: {ResistancePerMeter: 0.000161, Ampacity: 445},  // 4/0
}

// LookupAWG returns data for a gauge, and whether it was found.
func LookupAWG(gauge int) (AWGData, bool) {
	d, ok := AWGTable[gauge]
	return d, ok
}
