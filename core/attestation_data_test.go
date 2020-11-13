package core

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func TestToCoreAttestationData(t *testing.T) {
	type args struct {
		req *pb.SignBeaconAttestationRequest
	}
	tests := []struct {
		name string
		args args
		want *BeaconAttestation
	}{
		{
			name: "all filled",
			args: args{
				req: &pb.SignBeaconAttestationRequest{
					Id: &pb.SignBeaconAttestationRequest_PublicKey{
						PublicKey: []byte("mypublickey"),
					},
					Domain: []byte("domain"),
					Data: &pb.AttestationData{
						Slot:            12,
						CommitteeIndex:  12,
						BeaconBlockRoot: []byte("BeaconBlockRoot"),
						Source: &pb.Checkpoint{
							Epoch: 11,
							Root:  []byte("Root1"),
						},
						Target: &pb.Checkpoint{
							Epoch: 12,
							Root:  []byte("Root2"),
						},
					},
				},
			},
			want: &BeaconAttestation{
				Slot:            12,
				CommitteeIndex:  12,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &Checkpoint{
					Epoch: 11,
					Root:  []byte("Root1"),
				},
				Target: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root2"),
				},
			},
		},
		{
			name: "empty target",
			args: args{
				req: &pb.SignBeaconAttestationRequest{
					Id: &pb.SignBeaconAttestationRequest_PublicKey{
						PublicKey: []byte("mypublickey"),
					},
					Domain: []byte("domain"),
					Data: &pb.AttestationData{
						Slot:            12,
						CommitteeIndex:  12,
						BeaconBlockRoot: []byte("BeaconBlockRoot"),
						Source: &pb.Checkpoint{
							Epoch: 11,
							Root:  []byte("Root1"),
						},
					},
				},
			},
			want: &BeaconAttestation{
				Slot:            12,
				CommitteeIndex:  12,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &Checkpoint{
					Epoch: 11,
					Root:  []byte("Root1"),
				},
			},
		},
		{
			name: "empty source",
			args: args{
				req: &pb.SignBeaconAttestationRequest{
					Id: &pb.SignBeaconAttestationRequest_PublicKey{
						PublicKey: []byte("mypublickey"),
					},
					Domain: []byte("domain"),
					Data: &pb.AttestationData{
						Slot:            12,
						CommitteeIndex:  12,
						BeaconBlockRoot: []byte("BeaconBlockRoot"),
						Target: &pb.Checkpoint{
							Epoch: 11,
							Root:  []byte("Root1"),
						},
					},
				},
			},
			want: &BeaconAttestation{
				Slot:            12,
				CommitteeIndex:  12,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Target: &Checkpoint{
					Epoch: 11,
					Root:  []byte("Root1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ToCoreAttestationData(tt.args.req))
		})
	}
}

func TestCheckpoint_compare(t *testing.T) {
	type fields struct {
		Epoch uint64
		Root  []byte
	}
	type args struct {
		other *Checkpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "all equal",
			fields: fields{
				Epoch: 12,
				Root:  []byte("Root"),
			},
			args: args{
				other: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root"),
				},
			},
			want: true,
		},
		{
			name: "different epochs",
			fields: fields{
				Epoch: 11,
				Root:  []byte("Root"),
			},
			args: args{
				other: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root"),
				},
			},
			want: false,
		},
		{
			name: "different roots",
			fields: fields{
				Epoch: 12,
				Root:  []byte("Root1"),
			},
			args: args{
				other: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root2"),
				},
			},
			want: false,
		},
		{
			name: "all different",
			fields: fields{
				Epoch: 11,
				Root:  []byte("Root1"),
			},
			args: args{
				other: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root2"),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkpoint := &Checkpoint{
				Epoch: tt.fields.Epoch,
				Root:  tt.fields.Root,
			}
			require.Equal(t, tt.want, checkpoint.compare(tt.args.other))
		})
	}
}

func TestBeaconAttestation_Compare(t *testing.T) {
	type fields struct {
		Slot            uint64
		CommitteeIndex  uint64
		BeaconBlockRoot []byte
		Source          *Checkpoint
		Target          *Checkpoint
	}
	type args struct {
		other *BeaconAttestation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "equal",
			fields: fields{
				Slot:            12,
				CommitteeIndex:  12,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root1"),
				},
				Target: &Checkpoint{
					Epoch: 11,
					Root:  []byte("Root2"),
				},
			},
			args: args{
				other: &BeaconAttestation{
					Slot:            12,
					CommitteeIndex:  12,
					BeaconBlockRoot: []byte("BeaconBlockRoot"),
					Source: &Checkpoint{
						Epoch: 12,
						Root:  []byte("Root1"),
					},
					Target: &Checkpoint{
						Epoch: 11,
						Root:  []byte("Root2"),
					},
				},
			},
			want: true,
		},
		{
			name: "different slots",
			fields: fields{
				Slot:            1,
				CommitteeIndex:  12,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root1"),
				},
				Target: &Checkpoint{
					Epoch: 11,
					Root:  []byte("Root2"),
				},
			},
			args: args{
				other: &BeaconAttestation{
					Slot:            2,
					CommitteeIndex:  12,
					BeaconBlockRoot: []byte("BeaconBlockRoot"),
					Source: &Checkpoint{
						Epoch: 12,
						Root:  []byte("Root1"),
					},
					Target: &Checkpoint{
						Epoch: 11,
						Root:  []byte("Root2"),
					},
				},
			},
			want: false,
		},
		{
			name: "different committee index",
			fields: fields{
				Slot:            12,
				CommitteeIndex:  1,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root1"),
				},
				Target: &Checkpoint{
					Epoch: 11,
					Root:  []byte("Root2"),
				},
			},
			args: args{
				other: &BeaconAttestation{
					Slot:            12,
					CommitteeIndex:  2,
					BeaconBlockRoot: []byte("BeaconBlockRoot"),
					Source: &Checkpoint{
						Epoch: 12,
						Root:  []byte("Root1"),
					},
					Target: &Checkpoint{
						Epoch: 11,
						Root:  []byte("Root2"),
					},
				},
			},
			want: false,
		},
		{
			name: "different beacon block root",
			fields: fields{
				Slot:            12,
				CommitteeIndex:  12,
				BeaconBlockRoot: []byte("BeaconBlockRoot1"),
				Source: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root1"),
				},
				Target: &Checkpoint{
					Epoch: 11,
					Root:  []byte("Root2"),
				},
			},
			args: args{
				other: &BeaconAttestation{
					Slot:            12,
					CommitteeIndex:  12,
					BeaconBlockRoot: []byte("BeaconBlockRoot2"),
					Source: &Checkpoint{
						Epoch: 12,
						Root:  []byte("Root1"),
					},
					Target: &Checkpoint{
						Epoch: 11,
						Root:  []byte("Root2"),
					},
				},
			},
			want: false,
		},
		{
			name: "different source",
			fields: fields{
				Slot:            12,
				CommitteeIndex:  12,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &Checkpoint{
					Epoch: 1,
					Root:  []byte("Root1"),
				},
				Target: &Checkpoint{
					Epoch: 11,
					Root:  []byte("Root2"),
				},
			},
			args: args{
				other: &BeaconAttestation{
					Slot:            12,
					CommitteeIndex:  12,
					BeaconBlockRoot: []byte("BeaconBlockRoot"),
					Source: &Checkpoint{
						Epoch: 2,
						Root:  []byte("Root1"),
					},
					Target: &Checkpoint{
						Epoch: 11,
						Root:  []byte("Root2"),
					},
				},
			},
			want: false,
		},
		{
			name: "different target",
			fields: fields{
				Slot:            12,
				CommitteeIndex:  12,
				BeaconBlockRoot: []byte("BeaconBlockRoot"),
				Source: &Checkpoint{
					Epoch: 12,
					Root:  []byte("Root1"),
				},
				Target: &Checkpoint{
					Epoch: 1,
					Root:  []byte("Root2"),
				},
			},
			args: args{
				other: &BeaconAttestation{
					Slot:            12,
					CommitteeIndex:  12,
					BeaconBlockRoot: []byte("BeaconBlockRoot"),
					Source: &Checkpoint{
						Epoch: 12,
						Root:  []byte("Root1"),
					},
					Target: &Checkpoint{
						Epoch: 2,
						Root:  []byte("Root2"),
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			att := &BeaconAttestation{
				Slot:            tt.fields.Slot,
				CommitteeIndex:  tt.fields.CommitteeIndex,
				BeaconBlockRoot: tt.fields.BeaconBlockRoot,
				Source:          tt.fields.Source,
				Target:          tt.fields.Target,
			}
			if got := att.Compare(tt.args.other); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
