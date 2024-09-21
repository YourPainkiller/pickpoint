package usecase

import (
	"homework1/internal/dto"
	"homework1/internal/usecase/mock"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
)

func TestAcceptWithEmptyBase(t *testing.T) {
	type args struct {
		req  dto.AcceptOrderRequest
		resp dto.ListOrdersDto
	}

	tests := []struct {
		name    string
		args    args
		setup   func() (*mock.OrderRepositoryMock, *OrderUseCase)
		wantErr bool
	}{
		{
			name: "SuccesfullAcceptBox",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{Orders: []dto.OrderDto{}}, nil)
				repoMock.InsertOrdersMock.Expect(&dto.ListOrdersDto{
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
				}).Return(nil)
				return repoMock, oc
			},
			wantErr: false,
		},
		{
			name: "SuccesfullAcceptBagWithStretch",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{Orders: []dto.OrderDto{}}, nil)
				repoMock.InsertOrdersMock.Expect(&dto.ListOrdersDto{
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
				}).Return(nil)
				return repoMock, oc
			},
			wantErr: false,
		},
		{
			name: "SuccesfullAcceptStretch",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{Orders: []dto.OrderDto{}}, nil)
				repoMock.InsertOrdersMock.Expect(&dto.ListOrdersDto{
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
				}).Return(nil)
				return repoMock, oc
			},
			wantErr: false,
		},
		{
			name: "WrongDateAccept",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{Orders: []dto.OrderDto{}}, nil)
				return repoMock, oc
			},
			wantErr: true,
		},
		{
			name: "DoubleStretch",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{Orders: []dto.OrderDto{}}, nil)
				return repoMock, oc
			},
			wantErr: true,
		},
	}

	for _, itt := range tests {
		tt := itt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, oc := tt.setup()
			if err := oc.Accept(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAcceptReturn(t *testing.T) {
	type args struct {
		req    dto.AcceptReturnOrderRequest
		resp   dto.ListOrdersDto
		insert dto.ListOrdersDto
	}

	tests := []struct {
		name    string
		args    args
		setup   func() (*mock.OrderRepositoryMock, *OrderUseCase)
		wantErr bool
	}{
		{
			name: "SuccesfullAcceptReturn",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				repoMock.InsertOrdersMock.Expect(&dto.ListOrdersDto{
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
				}).Return(nil)

				return repoMock, oc

			},
			wantErr: false,
		},
		{
			name: "WrongOrderId",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				return repoMock, oc

			},
			wantErr: true,
		},
		{
			name: "WrongUserId",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				return repoMock, oc

			},
			wantErr: true,
		},
		{
			name: "BadTime",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				return repoMock, oc

			},
			wantErr: true,
		},
	}

	for _, tti := range tests {
		tt := tti
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, oc := tt.setup()
			if err := oc.AcceptReturn(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("AcceptReturn() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestGive(t *testing.T) {
	type args struct {
		req    dto.GiveOrderRequest
		resp   dto.ListOrdersDto
		insert dto.ListOrdersDto
	}

	tests := []struct {
		name    string
		args    args
		setup   func() (*mock.OrderRepositoryMock, *OrderUseCase)
		wantErr bool
	}{
		{
			name: "SuccesfullGiveOneOrder",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				repoMock.InsertOrdersMock.Expect(&dto.ListOrdersDto{
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
				}).Return(nil)
				return repoMock, oc
			},
			wantErr: false,
		},
		{
			name: "SuccesfullGiveMultipleOrders",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				repoMock.InsertOrdersMock.Expect(&dto.ListOrdersDto{
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
				}).Return(nil)
				return repoMock, oc
			},
			wantErr: false,
		},
		{
			name: "GiveOrdersToDifferentUsers",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				return repoMock, oc
			},
			wantErr: true,
		},
		{
			name: "WrongTime",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				return repoMock, oc
			},
			wantErr: true,
		},
		{
			name: "OrderGivedOrReturned",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				return repoMock, oc
			},
			wantErr: true,
		},
	}

	for _, tti := range tests {
		tt := tti
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, oc := tt.setup()
			if err := oc.Give(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Give() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestReturn(t *testing.T) {
	type args struct {
		req    dto.ReturnOrderRequest
		resp   dto.ListOrdersDto
		insert dto.ListOrdersDto
	}

	tests := []struct {
		name    string
		args    args
		setup   func() (*mock.OrderRepositoryMock, *OrderUseCase)
		wantErr bool
	}{
		{
			name: "SuccesfullReturnToCourier",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				repoMock.InsertOrdersMock.Expect(&dto.ListOrdersDto{
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
				}).Return(nil)
				return repoMock, oc
			},
			wantErr: false,
		},
		{
			name: "CantReturnOrderWhenGived",
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
			setup: func() (*mock.OrderRepositoryMock, *OrderUseCase) {
				ctrl := minimock.NewController(t)
				repoMock := mock.NewOrderRepositoryMock(ctrl)
				oc := &OrderUseCase{
					repo: repoMock,
				}
				repoMock.GetOrdersMock.Expect().Return(&dto.ListOrdersDto{
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
				}, nil)
				return repoMock, oc
			},
			wantErr: true,
		},
	}

	for _, tti := range tests {
		tt := tti
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, oc := tt.setup()
			if err := oc.Return(&tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Return() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
