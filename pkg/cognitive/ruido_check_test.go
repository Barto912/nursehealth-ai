package cognitive

import "testing"

func TestAllowsCareEvent(t *testing.T) {
	cases := []struct {
		name string
		rc   RuidoCheck
		want bool
	}{
		{"clear and moderate passes", RuidoCheck{MentalClarity: 8, WorkLevel: WorkLevelModerate}, true},
		{"boundary clarity 7 passes", RuidoCheck{MentalClarity: 7, WorkLevel: WorkLevelHigh}, true},
		{"clarity 6 blocks", RuidoCheck{MentalClarity: 6, WorkLevel: WorkLevelHigh}, false},
		{"low work level blocks", RuidoCheck{MentalClarity: 9, WorkLevel: WorkLevelLow}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.rc.AllowsCareEvent(); got != tc.want {
				t.Fatalf("AllowsCareEvent() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAllowsSNOMEDSignature(t *testing.T) {
	cases := []struct {
		name string
		rc   RuidoCheck
		want bool
	}{
		{"simple attention passes", RuidoCheck{SimpleAttention: true}, true},
		{"no attention blocks", RuidoCheck{SimpleAttention: false}, false},
		{"finding rescues signature", RuidoCheck{SimpleAttention: false, AttentionalFinding: "Paciente consciente, orientado, acompaña familiar"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.rc.AllowsSNOMEDSignature(); got != tc.want {
				t.Fatalf("AllowsSNOMEDSignature() = %v, want %v", got, tc.want)
			}
		})
	}
}