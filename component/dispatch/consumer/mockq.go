// Author: XinRui Hua
// Time:   2022/1/27 下午3:05
// Git:    huaxr

package consumer

type mockConsumer struct {
}

func (h *mockConsumer) Consume() {

}

func (h *mockConsumer) Run() error {
	return nil
}

func (h *mockConsumer) SetAuth(u string, p string) {}

func NewMockConsumer() *mockConsumer {
	return nil
}
