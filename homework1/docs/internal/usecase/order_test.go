package usecase

import (
	"homework1/internal/dto"
	"homework1/internal/usecase/mock"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
)

func TestAcceptWithEmptyBase(t *testing.T) {
	type fields struct {
		repo orderRepository
	}

	type args struct {
		req  dto.AcceptOrderRequest
		resp dto.ListOrdersDto
	}

	ctrl := minimock.NewController(t)
	repoMock := mock.NewOrderRepositoryMock(ctrl)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "SuccesfullAcceptBox",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptOrderRequest{
					Id:                1,
					UserId:            1,
					ValidTime:         "2024-09-21",
					Price:             100,
					Weight:            100,
					PackageType:       "box",
					AdditionalStretch: false,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "SuccesfullAcceptBagWithStretch",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptOrderRequest{
					Id:                1,
					UserId:            1,
					ValidTime:         "2024-09-21",
					Price:             100,
					Weight:            100,
					PackageType:       "bag",
					AdditionalStretch: true,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             106,
							Weight:            100,
							PackageType:       "bag",
							AdditionalStretch: true,
							State:             "accepted",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "SuccesfullAcceptStretch",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptOrderRequest{
					Id:                1,
					UserId:            1,
					ValidTime:         "2024-09-21",
					Price:             100,
					Weight:            100,
					PackageType:       "stretch",
					AdditionalStretch: false,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             101,
							Weight:            100,
							PackageType:       "stretch",
							AdditionalStretch: false,
							State:             "accepted",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "WrongDateAccept",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptOrderRequest{
					Id:                1,
					UserId:            1,
					ValidTime:         "2000-09-21",
					Price:             100,
					Weight:            100,
					PackageType:       "stretch",
					AdditionalStretch: false,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             101,
							Weight:            100,
							PackageType:       "stretch",
							AdditionalStretch: false,
							State:             "accepted",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "DoubleStretch",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptOrderRequest{
					Id:                1,
					UserId:            1,
					ValidTime:         "2000-09-21",
					Price:             100,
					Weight:            100,
					PackageType:       "stretch",
					AdditionalStretch: true,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             101,
							Weight:            100,
							PackageType:       "stretch",
							AdditionalStretch: true,
							State:             "accepted",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc := &OrderUseCase{
				repo: tt.fields.repo,
			}
			repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{Orders: []dto.OrderDto{}}, nil)
			repoMock.InsertOrdersMock.Expect(&tt.args.resp).Return(nil)

			if err := oc.Accept(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAcceptWithOrderInBase(t *testing.T) {
	type fields struct {
		repo orderRepository
	}

	type args struct {
		req  dto.AcceptOrderRequest
		resp dto.ListOrdersDto
	}

	ctrl := minimock.NewController(t)
	repoMock := mock.NewOrderRepositoryMock(ctrl)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "SameOrderError",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptOrderRequest{
					Id:                2,
					UserId:            1,
					ValidTime:         "2024-09-21",
					Price:             100,
					Weight:            100,
					PackageType:       "box",
					AdditionalStretch: false,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                2,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc := &OrderUseCase{
				repo: tt.fields.repo,
			}
			repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{Orders: []dto.OrderDto{{Id: 2}}}, nil)
			//repoMock.InsertOrdersMock.Expect(&tt.args.resp).Return(nil)

			if err := oc.Accept(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAcceptReturn(t *testing.T) {
	type fields struct {
		repo orderRepository
	}

	type args struct {
		req    dto.AcceptReturnOrderRequest
		resp   dto.ListOrdersDto
		insert dto.ListOrdersDto
	}

	ctrl := minimock.NewController(t)
	repoMock := mock.NewOrderRepositoryMock(ctrl)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "SuccesfullAcceptReturn",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptReturnOrderRequest{
					Id:     1,
					UserId: 1,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "WrongOrderId",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptReturnOrderRequest{
					Id:     1,
					UserId: 1,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                2,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "WrongUserId",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptReturnOrderRequest{
					Id:     2,
					UserId: 1,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "BadTime",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.AcceptReturnOrderRequest{
					Id:     1,
					UserId: 1,
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2000-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc := &OrderUseCase{
				repo: tt.fields.repo,
			}
			repoMock.GetOrdersMock.Expect().Return(&tt.args.resp, nil)

			repoMock.InsertOrdersMock.Expect(&tt.args.insert).Return(nil)

			if err := oc.AcceptReturn(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestGive(t *testing.T) {
	type fields struct {
		repo orderRepository
	}

	type args struct {
		req    dto.GiveOrderRequest
		resp   dto.ListOrdersDto
		insert dto.ListOrdersDto
	}

	ctrl := minimock.NewController(t)
	repoMock := mock.NewOrderRepositoryMock(ctrl)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "SuccesfullGiveOneOrder",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.GiveOrderRequest{
					OrderIds: []dto.OrderId{
						{
							Id: 1,
						},
					},
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "SuccesfullGiveMultipleOrders",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.GiveOrderRequest{
					OrderIds: []dto.OrderId{
						{
							Id: 1,
						},
						{
							Id: 2,
						},
					},
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
						{
							Id:                2,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
						{
							Id:                2,
							UserId:            1,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "GiveOrdersToDifferentUsers",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.GiveOrderRequest{
					OrderIds: []dto.OrderId{
						{
							Id: 1,
						},
						{
							Id: 2,
						},
					},
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
						{
							Id:                2,
							UserId:            2,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
						{
							Id:                2,
							UserId:            2,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "WrongTime",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.GiveOrderRequest{
					OrderIds: []dto.OrderId{
						{
							Id: 1,
						},
						{
							Id: 2,
						},
					},
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2000-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
						{
							Id:                2,
							UserId:            1,
							ValidTime:         "2000-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "accepted",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
						{
							Id:                2,
							UserId:            1,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name:   "OrderGivedOrReturned",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.GiveOrderRequest{
					OrderIds: []dto.OrderId{
						{
							Id: 1,
						},
						{
							Id: 2,
						},
					},
				},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2000-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
						{
							Id:                2,
							UserId:            1,
							ValidTime:         "2000-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
						{
							Id:                2,
							UserId:            1,
							ValidTime:         time.Now().Add(time.Hour * 72).Format(TIMELAYOUT),
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc := &OrderUseCase{
				repo: tt.fields.repo,
			}
			repoMock.GetOrdersMock.Expect().Return(&tt.args.resp, nil)

			repoMock.InsertOrdersMock.Expect(&tt.args.insert).Return(nil)

			if err := oc.Give(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestReturn(t *testing.T) {
	type fields struct {
		repo orderRepository
	}

	type args struct {
		req    dto.ReturnOrderRequest
		resp   dto.ListOrdersDto
		insert dto.ListOrdersDto
	}

	ctrl := minimock.NewController(t)
	repoMock := mock.NewOrderRepositoryMock(ctrl)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "SuccesfullReturnToCourier",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.ReturnOrderRequest{Id: 1},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
						{
							Id:                2,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                2,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "CantReturnOrderWhenGived",
			fields: fields{repo: repoMock},
			args: args{
				req: dto.ReturnOrderRequest{Id: 1},
				resp: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{
							Id:                1,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "gived",
						},
						{
							Id:                2,
							UserId:            1,
							ValidTime:         "2024-09-22",
							Price:             120,
							Weight:            100,
							PackageType:       "box",
							AdditionalStretch: false,
							State:             "returned",
						},
					},
				},
				insert: dto.ListOrdersDto{
					Orders: []dto.OrderDto{
						{},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oc := &OrderUseCase{
				repo: tt.fields.repo,
			}
			repoMock.GetOrdersMock.Expect().Return(&tt.args.resp, nil)

			repoMock.InsertOrdersMock.Expect(&tt.args.insert).Return(nil)

			if err := oc.Return(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
