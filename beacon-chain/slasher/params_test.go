package slasher

import (
	"reflect"
	"testing"

	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
)

func TestDefaultParams(t *testing.T) {
	def := DefaultParams()
	assert.Equal(t, true, def.chunkSize > 0)
	assert.Equal(t, true, def.validatorChunkSize > 0)
	assert.Equal(t, true, def.historyLength > 0)
}

func TestParams_cellIndex(t *testing.T) {
	type args struct {
		validatorIndex types.ValidatorIndex
		epoch          types.Epoch
	}
	tests := []struct {
		name   string
		fields *Parameters
		args   args
		want   uint64
	}{
		{
			name: "epoch 0 and validator index 0",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
			},
			args: args{
				validatorIndex: 0,
				epoch:          0,
			},
			want: 0,
		},
		{
			//     val0     val1     val2
			//      |        |        |
			//   {     }  {     }  {     }
			//  [2, 2, 2, 2, 2, 2, 2, 2, 2]
			//                        |-> epoch 1, validator 2
			name: "epoch < chunkSize and validatorIndex < validatorChunkSize",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
			},
			args: args{
				validatorIndex: 2,
				epoch:          1,
			},
			want: 7,
		},
		{
			//     val0     val1     val2
			//      |        |        |
			//   {     }  {     }  {     }
			//  [2, 2, 2, 2, 2, 2, 2, 2, 2]
			//                        |-> epoch 4, validator 2 (wrap around)
			name: "epoch > chunkSize and validatorIndex < validatorChunkSize",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
			},
			args: args{
				validatorIndex: 2,
				epoch:          4,
			},
			want: 7,
		},
		{
			//     val0     val1     val2
			//      |        |        |
			//   {     }  {     }  {     }
			//  [2, 2, 2, 2, 2, 2, 2, 2, 2]
			//                     |-> epoch 3, validator 2 (wrap around)
			name: "epoch = chunkSize and validatorIndex < validatorChunkSize",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
			},
			args: args{
				validatorIndex: 2,
				epoch:          3,
			},
			want: 6,
		},
		{
			//     val0     val1     val2
			//      |        |        |
			//   {     }  {     }  {     }
			//  [2, 2, 2, 2, 2, 2, 2, 2, 2]
			//   |-> epoch 0, validator 3 (wrap around)
			name: "epoch < chunkSize and validatorIndex = validatorChunkSize",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
			},
			args: args{
				validatorIndex: 3,
				epoch:          0,
			},
			want: 0,
		},
		{
			//     val0     val1     val2
			//      |        |        |
			//   {     }  {     }  {     }
			//  [2, 2, 2, 2, 2, 2, 2, 2, 2]
			//            |-> epoch 0, validator 4 (wrap around)
			name: "epoch < chunkSize and validatorIndex > validatorChunkSize",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
			},
			args: args{
				validatorIndex: 4,
				epoch:          0,
			},
			want: 3,
		},
		{
			//     val0     val1     val2
			//      |        |        |
			//   {     }  {     }  {     }
			//  [2, 2, 2, 2, 2, 2, 2, 2, 2]
			//   |-> epoch 3, validator 3 (wrap around)
			name: "epoch = chunkSize and validatorIndex = validatorChunkSize",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
			},
			args: args{
				validatorIndex: 3,
				epoch:          3,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Parameters{
				chunkSize:          tt.fields.chunkSize,
				validatorChunkSize: tt.fields.validatorChunkSize,
				historyLength:      tt.fields.historyLength,
			}
			if got := c.cellIndex(tt.args.validatorIndex, tt.args.epoch); got != tt.want {
				t.Errorf("cellIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_chunkIndex(t *testing.T) {
	tests := []struct {
		name   string
		fields *Parameters
		epoch  types.Epoch
		want   uint64
	}{
		{
			name: "epoch 0",
			fields: &Parameters{
				chunkSize:     3,
				historyLength: 3,
			},
			epoch: 0,
			want:  0,
		},
		{
			name: "epoch < historyLength, epoch < chunkSize",
			fields: &Parameters{
				chunkSize:     3,
				historyLength: 3,
			},
			epoch: 2,
			want:  0,
		},
		{
			name: "epoch = historyLength, epoch < chunkSize",
			fields: &Parameters{
				chunkSize:     4,
				historyLength: 3,
			},
			epoch: 3,
			want:  0,
		},
		{
			name: "epoch > historyLength, epoch < chunkSize",
			fields: &Parameters{
				chunkSize:     5,
				historyLength: 3,
			},
			epoch: 4,
			want:  0,
		},
		{
			name: "epoch < historyLength, epoch < chunkSize",
			fields: &Parameters{
				chunkSize:     3,
				historyLength: 3,
			},
			epoch: 2,
			want:  0,
		},
		{
			name: "epoch = historyLength, epoch < chunkSize",
			fields: &Parameters{
				chunkSize:     4,
				historyLength: 3,
			},
			epoch: 3,
			want:  0,
		},
		{
			name: "epoch < historyLength, epoch = chunkSize",
			fields: &Parameters{
				chunkSize:     2,
				historyLength: 3,
			},
			epoch: 2,
			want:  1,
		},
		{
			name: "epoch < historyLength, epoch > chunkSize",
			fields: &Parameters{
				chunkSize:     2,
				historyLength: 4,
			},
			epoch: 3,
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Parameters{
				chunkSize:     tt.fields.chunkSize,
				historyLength: tt.fields.historyLength,
			}
			if got := c.chunkIndex(tt.epoch); got != tt.want {
				t.Errorf("chunkIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_flatSliceID(t *testing.T) {
	tests := []struct {
		name                string
		fields              *Parameters
		validatorChunkIndex uint64
		chunkIndex          uint64
		want                uint64
	}{
		{
			name: "Proper disk key for 0, 0",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
				historyLength:      6,
			},
			chunkIndex:          0,
			validatorChunkIndex: 0,
			want:                0,
		},
		{
			name: "Proper disk key for epoch < historyLength, validator < validatorChunkSize",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
				historyLength:      6,
			},
			chunkIndex:          1,
			validatorChunkIndex: 1,
			want:                3,
		},
		{
			name: "Proper disk key for epoch > historyLength, validator > validatorChunkSize",
			fields: &Parameters{
				chunkSize:          3,
				validatorChunkSize: 3,
				historyLength:      6,
			},
			chunkIndex:          10,
			validatorChunkIndex: 10,
			want:                30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Parameters{
				chunkSize:          tt.fields.chunkSize,
				validatorChunkSize: tt.fields.validatorChunkSize,
				historyLength:      tt.fields.historyLength,
			}
			if got := c.flatSliceID(tt.validatorChunkIndex, tt.chunkIndex); got != tt.want {
				t.Errorf("diskKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_validatorChunkIndex(t *testing.T) {
	tests := []struct {
		name           string
		fields         *Parameters
		validatorIndex types.ValidatorIndex
		want           uint64
	}{
		{
			name: "validator index < validatorChunkSize",
			fields: &Parameters{
				validatorChunkSize: 3,
			},
			validatorIndex: 2,
			want:           0,
		},
		{
			name: "validator index = validatorChunkSize",
			fields: &Parameters{
				validatorChunkSize: 3,
			},
			validatorIndex: 3,
			want:           1,
		},
		{
			name: "validator index > validatorChunkSize",
			fields: &Parameters{
				validatorChunkSize: 3,
			},
			validatorIndex: 99,
			want:           33,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Parameters{
				validatorChunkSize: tt.fields.validatorChunkSize,
			}
			if got := c.validatorChunkIndex(tt.validatorIndex); got != tt.want {
				t.Errorf("validatorChunkIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_chunkOffset(t *testing.T) {
	tests := []struct {
		name   string
		fields *Parameters
		epoch  types.Epoch
		want   uint64
	}{
		{
			name: "epoch < chunkSize",
			fields: &Parameters{
				chunkSize: 3,
			},
			epoch: 2,
			want:  2,
		},
		{
			name: "epoch = chunkSize",
			fields: &Parameters{
				chunkSize: 3,
			},
			epoch: 3,
			want:  0,
		},
		{
			name: "epoch > chunkSize",
			fields: &Parameters{
				chunkSize: 3,
			},
			epoch: 5,
			want:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Parameters{
				chunkSize: tt.fields.chunkSize,
			}
			if got := c.chunkOffset(tt.epoch); got != tt.want {
				t.Errorf("chunkOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_validatorOffset(t *testing.T) {
	tests := []struct {
		name           string
		fields         *Parameters
		validatorIndex types.ValidatorIndex
		want           uint64
	}{
		{
			name: "validatorIndex < validatorChunkSize",
			fields: &Parameters{
				validatorChunkSize: 3,
			},
			validatorIndex: 2,
			want:           2,
		},
		{
			name: "validatorIndex = validatorChunkSize",
			fields: &Parameters{
				validatorChunkSize: 3,
			},
			validatorIndex: 3,
			want:           0,
		},
		{
			name: "validatorIndex > validatorChunkSize",
			fields: &Parameters{
				validatorChunkSize: 3,
			},
			validatorIndex: 5,
			want:           2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Parameters{
				validatorChunkSize: tt.fields.validatorChunkSize,
			}
			if got := c.validatorOffset(tt.validatorIndex); got != tt.want {
				t.Errorf("validatorOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_validatorIndicesInChunk(t *testing.T) {
	tests := []struct {
		name              string
		fields            *Parameters
		validatorChunkIdx uint64
		want              []types.ValidatorIndex
	}{
		{
			name: "Returns proper indices",
			fields: &Parameters{
				validatorChunkSize: 3,
			},
			validatorChunkIdx: 2,
			want:              []types.ValidatorIndex{6, 7, 8},
		},
		{
			name: "0 validator chunk size returs empty",
			fields: &Parameters{
				validatorChunkSize: 0,
			},
			validatorChunkIdx: 100,
			want:              []types.ValidatorIndex{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Parameters{
				validatorChunkSize: tt.fields.validatorChunkSize,
			}
			if got := c.validatorIndicesInChunk(tt.validatorChunkIdx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validatorIndicesInChunk() = %v, want %v", got, tt.want)
			}
		})
	}
}
