// Code generated by http://github.com/gojuno/minimock (v3.4.0). DO NOT EDIT.

package mocks

//go:generate minimock -i gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/service.LomsClient -o loms_client_mock.go -n LomsClientMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
)

// LomsClientMock implements mm_service.LomsClient
type LomsClientMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetStock          func(ctx context.Context, skuID int64) (i1 int64, err error)
	funcGetStockOrigin    string
	inspectFuncGetStock   func(ctx context.Context, skuID int64)
	afterGetStockCounter  uint64
	beforeGetStockCounter uint64
	GetStockMock          mLomsClientMockGetStock

	funcOrderCreate          func(ctx context.Context, userID int64, items []model.CartItem) (orderID int64, err error)
	funcOrderCreateOrigin    string
	inspectFuncOrderCreate   func(ctx context.Context, userID int64, items []model.CartItem)
	afterOrderCreateCounter  uint64
	beforeOrderCreateCounter uint64
	OrderCreateMock          mLomsClientMockOrderCreate
}

// NewLomsClientMock returns a mock for mm_service.LomsClient
func NewLomsClientMock(t minimock.Tester) *LomsClientMock {
	m := &LomsClientMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetStockMock = mLomsClientMockGetStock{mock: m}
	m.GetStockMock.callArgs = []*LomsClientMockGetStockParams{}

	m.OrderCreateMock = mLomsClientMockOrderCreate{mock: m}
	m.OrderCreateMock.callArgs = []*LomsClientMockOrderCreateParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mLomsClientMockGetStock struct {
	optional           bool
	mock               *LomsClientMock
	defaultExpectation *LomsClientMockGetStockExpectation
	expectations       []*LomsClientMockGetStockExpectation

	callArgs []*LomsClientMockGetStockParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// LomsClientMockGetStockExpectation specifies expectation struct of the LomsClient.GetStock
type LomsClientMockGetStockExpectation struct {
	mock               *LomsClientMock
	params             *LomsClientMockGetStockParams
	paramPtrs          *LomsClientMockGetStockParamPtrs
	expectationOrigins LomsClientMockGetStockExpectationOrigins
	results            *LomsClientMockGetStockResults
	returnOrigin       string
	Counter            uint64
}

// LomsClientMockGetStockParams contains parameters of the LomsClient.GetStock
type LomsClientMockGetStockParams struct {
	ctx   context.Context
	skuID int64
}

// LomsClientMockGetStockParamPtrs contains pointers to parameters of the LomsClient.GetStock
type LomsClientMockGetStockParamPtrs struct {
	ctx   *context.Context
	skuID *int64
}

// LomsClientMockGetStockResults contains results of the LomsClient.GetStock
type LomsClientMockGetStockResults struct {
	i1  int64
	err error
}

// LomsClientMockGetStockOrigins contains origins of expectations of the LomsClient.GetStock
type LomsClientMockGetStockExpectationOrigins struct {
	origin      string
	originCtx   string
	originSkuID string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmGetStock *mLomsClientMockGetStock) Optional() *mLomsClientMockGetStock {
	mmGetStock.optional = true
	return mmGetStock
}

// Expect sets up expected params for LomsClient.GetStock
func (mmGetStock *mLomsClientMockGetStock) Expect(ctx context.Context, skuID int64) *mLomsClientMockGetStock {
	if mmGetStock.mock.funcGetStock != nil {
		mmGetStock.mock.t.Fatalf("LomsClientMock.GetStock mock is already set by Set")
	}

	if mmGetStock.defaultExpectation == nil {
		mmGetStock.defaultExpectation = &LomsClientMockGetStockExpectation{}
	}

	if mmGetStock.defaultExpectation.paramPtrs != nil {
		mmGetStock.mock.t.Fatalf("LomsClientMock.GetStock mock is already set by ExpectParams functions")
	}

	mmGetStock.defaultExpectation.params = &LomsClientMockGetStockParams{ctx, skuID}
	mmGetStock.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmGetStock.expectations {
		if minimock.Equal(e.params, mmGetStock.defaultExpectation.params) {
			mmGetStock.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetStock.defaultExpectation.params)
		}
	}

	return mmGetStock
}

// ExpectCtxParam1 sets up expected param ctx for LomsClient.GetStock
func (mmGetStock *mLomsClientMockGetStock) ExpectCtxParam1(ctx context.Context) *mLomsClientMockGetStock {
	if mmGetStock.mock.funcGetStock != nil {
		mmGetStock.mock.t.Fatalf("LomsClientMock.GetStock mock is already set by Set")
	}

	if mmGetStock.defaultExpectation == nil {
		mmGetStock.defaultExpectation = &LomsClientMockGetStockExpectation{}
	}

	if mmGetStock.defaultExpectation.params != nil {
		mmGetStock.mock.t.Fatalf("LomsClientMock.GetStock mock is already set by Expect")
	}

	if mmGetStock.defaultExpectation.paramPtrs == nil {
		mmGetStock.defaultExpectation.paramPtrs = &LomsClientMockGetStockParamPtrs{}
	}
	mmGetStock.defaultExpectation.paramPtrs.ctx = &ctx
	mmGetStock.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmGetStock
}

// ExpectSkuIDParam2 sets up expected param skuID for LomsClient.GetStock
func (mmGetStock *mLomsClientMockGetStock) ExpectSkuIDParam2(skuID int64) *mLomsClientMockGetStock {
	if mmGetStock.mock.funcGetStock != nil {
		mmGetStock.mock.t.Fatalf("LomsClientMock.GetStock mock is already set by Set")
	}

	if mmGetStock.defaultExpectation == nil {
		mmGetStock.defaultExpectation = &LomsClientMockGetStockExpectation{}
	}

	if mmGetStock.defaultExpectation.params != nil {
		mmGetStock.mock.t.Fatalf("LomsClientMock.GetStock mock is already set by Expect")
	}

	if mmGetStock.defaultExpectation.paramPtrs == nil {
		mmGetStock.defaultExpectation.paramPtrs = &LomsClientMockGetStockParamPtrs{}
	}
	mmGetStock.defaultExpectation.paramPtrs.skuID = &skuID
	mmGetStock.defaultExpectation.expectationOrigins.originSkuID = minimock.CallerInfo(1)

	return mmGetStock
}

// Inspect accepts an inspector function that has same arguments as the LomsClient.GetStock
func (mmGetStock *mLomsClientMockGetStock) Inspect(f func(ctx context.Context, skuID int64)) *mLomsClientMockGetStock {
	if mmGetStock.mock.inspectFuncGetStock != nil {
		mmGetStock.mock.t.Fatalf("Inspect function is already set for LomsClientMock.GetStock")
	}

	mmGetStock.mock.inspectFuncGetStock = f

	return mmGetStock
}

// Return sets up results that will be returned by LomsClient.GetStock
func (mmGetStock *mLomsClientMockGetStock) Return(i1 int64, err error) *LomsClientMock {
	if mmGetStock.mock.funcGetStock != nil {
		mmGetStock.mock.t.Fatalf("LomsClientMock.GetStock mock is already set by Set")
	}

	if mmGetStock.defaultExpectation == nil {
		mmGetStock.defaultExpectation = &LomsClientMockGetStockExpectation{mock: mmGetStock.mock}
	}
	mmGetStock.defaultExpectation.results = &LomsClientMockGetStockResults{i1, err}
	mmGetStock.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmGetStock.mock
}

// Set uses given function f to mock the LomsClient.GetStock method
func (mmGetStock *mLomsClientMockGetStock) Set(f func(ctx context.Context, skuID int64) (i1 int64, err error)) *LomsClientMock {
	if mmGetStock.defaultExpectation != nil {
		mmGetStock.mock.t.Fatalf("Default expectation is already set for the LomsClient.GetStock method")
	}

	if len(mmGetStock.expectations) > 0 {
		mmGetStock.mock.t.Fatalf("Some expectations are already set for the LomsClient.GetStock method")
	}

	mmGetStock.mock.funcGetStock = f
	mmGetStock.mock.funcGetStockOrigin = minimock.CallerInfo(1)
	return mmGetStock.mock
}

// When sets expectation for the LomsClient.GetStock which will trigger the result defined by the following
// Then helper
func (mmGetStock *mLomsClientMockGetStock) When(ctx context.Context, skuID int64) *LomsClientMockGetStockExpectation {
	if mmGetStock.mock.funcGetStock != nil {
		mmGetStock.mock.t.Fatalf("LomsClientMock.GetStock mock is already set by Set")
	}

	expectation := &LomsClientMockGetStockExpectation{
		mock:               mmGetStock.mock,
		params:             &LomsClientMockGetStockParams{ctx, skuID},
		expectationOrigins: LomsClientMockGetStockExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmGetStock.expectations = append(mmGetStock.expectations, expectation)
	return expectation
}

// Then sets up LomsClient.GetStock return parameters for the expectation previously defined by the When method
func (e *LomsClientMockGetStockExpectation) Then(i1 int64, err error) *LomsClientMock {
	e.results = &LomsClientMockGetStockResults{i1, err}
	return e.mock
}

// Times sets number of times LomsClient.GetStock should be invoked
func (mmGetStock *mLomsClientMockGetStock) Times(n uint64) *mLomsClientMockGetStock {
	if n == 0 {
		mmGetStock.mock.t.Fatalf("Times of LomsClientMock.GetStock mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmGetStock.expectedInvocations, n)
	mmGetStock.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmGetStock
}

func (mmGetStock *mLomsClientMockGetStock) invocationsDone() bool {
	if len(mmGetStock.expectations) == 0 && mmGetStock.defaultExpectation == nil && mmGetStock.mock.funcGetStock == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmGetStock.mock.afterGetStockCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmGetStock.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// GetStock implements mm_service.LomsClient
func (mmGetStock *LomsClientMock) GetStock(ctx context.Context, skuID int64) (i1 int64, err error) {
	mm_atomic.AddUint64(&mmGetStock.beforeGetStockCounter, 1)
	defer mm_atomic.AddUint64(&mmGetStock.afterGetStockCounter, 1)

	mmGetStock.t.Helper()

	if mmGetStock.inspectFuncGetStock != nil {
		mmGetStock.inspectFuncGetStock(ctx, skuID)
	}

	mm_params := LomsClientMockGetStockParams{ctx, skuID}

	// Record call args
	mmGetStock.GetStockMock.mutex.Lock()
	mmGetStock.GetStockMock.callArgs = append(mmGetStock.GetStockMock.callArgs, &mm_params)
	mmGetStock.GetStockMock.mutex.Unlock()

	for _, e := range mmGetStock.GetStockMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.i1, e.results.err
		}
	}

	if mmGetStock.GetStockMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetStock.GetStockMock.defaultExpectation.Counter, 1)
		mm_want := mmGetStock.GetStockMock.defaultExpectation.params
		mm_want_ptrs := mmGetStock.GetStockMock.defaultExpectation.paramPtrs

		mm_got := LomsClientMockGetStockParams{ctx, skuID}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmGetStock.t.Errorf("LomsClientMock.GetStock got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmGetStock.GetStockMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.skuID != nil && !minimock.Equal(*mm_want_ptrs.skuID, mm_got.skuID) {
				mmGetStock.t.Errorf("LomsClientMock.GetStock got unexpected parameter skuID, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmGetStock.GetStockMock.defaultExpectation.expectationOrigins.originSkuID, *mm_want_ptrs.skuID, mm_got.skuID, minimock.Diff(*mm_want_ptrs.skuID, mm_got.skuID))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetStock.t.Errorf("LomsClientMock.GetStock got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmGetStock.GetStockMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetStock.GetStockMock.defaultExpectation.results
		if mm_results == nil {
			mmGetStock.t.Fatal("No results are set for the LomsClientMock.GetStock")
		}
		return (*mm_results).i1, (*mm_results).err
	}
	if mmGetStock.funcGetStock != nil {
		return mmGetStock.funcGetStock(ctx, skuID)
	}
	mmGetStock.t.Fatalf("Unexpected call to LomsClientMock.GetStock. %v %v", ctx, skuID)
	return
}

// GetStockAfterCounter returns a count of finished LomsClientMock.GetStock invocations
func (mmGetStock *LomsClientMock) GetStockAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetStock.afterGetStockCounter)
}

// GetStockBeforeCounter returns a count of LomsClientMock.GetStock invocations
func (mmGetStock *LomsClientMock) GetStockBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetStock.beforeGetStockCounter)
}

// Calls returns a list of arguments used in each call to LomsClientMock.GetStock.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetStock *mLomsClientMockGetStock) Calls() []*LomsClientMockGetStockParams {
	mmGetStock.mutex.RLock()

	argCopy := make([]*LomsClientMockGetStockParams, len(mmGetStock.callArgs))
	copy(argCopy, mmGetStock.callArgs)

	mmGetStock.mutex.RUnlock()

	return argCopy
}

// MinimockGetStockDone returns true if the count of the GetStock invocations corresponds
// the number of defined expectations
func (m *LomsClientMock) MinimockGetStockDone() bool {
	if m.GetStockMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.GetStockMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.GetStockMock.invocationsDone()
}

// MinimockGetStockInspect logs each unmet expectation
func (m *LomsClientMock) MinimockGetStockInspect() {
	for _, e := range m.GetStockMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LomsClientMock.GetStock at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterGetStockCounter := mm_atomic.LoadUint64(&m.afterGetStockCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.GetStockMock.defaultExpectation != nil && afterGetStockCounter < 1 {
		if m.GetStockMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to LomsClientMock.GetStock at\n%s", m.GetStockMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to LomsClientMock.GetStock at\n%s with params: %#v", m.GetStockMock.defaultExpectation.expectationOrigins.origin, *m.GetStockMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetStock != nil && afterGetStockCounter < 1 {
		m.t.Errorf("Expected call to LomsClientMock.GetStock at\n%s", m.funcGetStockOrigin)
	}

	if !m.GetStockMock.invocationsDone() && afterGetStockCounter > 0 {
		m.t.Errorf("Expected %d calls to LomsClientMock.GetStock at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.GetStockMock.expectedInvocations), m.GetStockMock.expectedInvocationsOrigin, afterGetStockCounter)
	}
}

type mLomsClientMockOrderCreate struct {
	optional           bool
	mock               *LomsClientMock
	defaultExpectation *LomsClientMockOrderCreateExpectation
	expectations       []*LomsClientMockOrderCreateExpectation

	callArgs []*LomsClientMockOrderCreateParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// LomsClientMockOrderCreateExpectation specifies expectation struct of the LomsClient.OrderCreate
type LomsClientMockOrderCreateExpectation struct {
	mock               *LomsClientMock
	params             *LomsClientMockOrderCreateParams
	paramPtrs          *LomsClientMockOrderCreateParamPtrs
	expectationOrigins LomsClientMockOrderCreateExpectationOrigins
	results            *LomsClientMockOrderCreateResults
	returnOrigin       string
	Counter            uint64
}

// LomsClientMockOrderCreateParams contains parameters of the LomsClient.OrderCreate
type LomsClientMockOrderCreateParams struct {
	ctx    context.Context
	userID int64
	items  []model.CartItem
}

// LomsClientMockOrderCreateParamPtrs contains pointers to parameters of the LomsClient.OrderCreate
type LomsClientMockOrderCreateParamPtrs struct {
	ctx    *context.Context
	userID *int64
	items  *[]model.CartItem
}

// LomsClientMockOrderCreateResults contains results of the LomsClient.OrderCreate
type LomsClientMockOrderCreateResults struct {
	orderID int64
	err     error
}

// LomsClientMockOrderCreateOrigins contains origins of expectations of the LomsClient.OrderCreate
type LomsClientMockOrderCreateExpectationOrigins struct {
	origin       string
	originCtx    string
	originUserID string
	originItems  string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmOrderCreate *mLomsClientMockOrderCreate) Optional() *mLomsClientMockOrderCreate {
	mmOrderCreate.optional = true
	return mmOrderCreate
}

// Expect sets up expected params for LomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) Expect(ctx context.Context, userID int64, items []model.CartItem) *mLomsClientMockOrderCreate {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &LomsClientMockOrderCreateExpectation{}
	}

	if mmOrderCreate.defaultExpectation.paramPtrs != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by ExpectParams functions")
	}

	mmOrderCreate.defaultExpectation.params = &LomsClientMockOrderCreateParams{ctx, userID, items}
	mmOrderCreate.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmOrderCreate.expectations {
		if minimock.Equal(e.params, mmOrderCreate.defaultExpectation.params) {
			mmOrderCreate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmOrderCreate.defaultExpectation.params)
		}
	}

	return mmOrderCreate
}

// ExpectCtxParam1 sets up expected param ctx for LomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) ExpectCtxParam1(ctx context.Context) *mLomsClientMockOrderCreate {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &LomsClientMockOrderCreateExpectation{}
	}

	if mmOrderCreate.defaultExpectation.params != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Expect")
	}

	if mmOrderCreate.defaultExpectation.paramPtrs == nil {
		mmOrderCreate.defaultExpectation.paramPtrs = &LomsClientMockOrderCreateParamPtrs{}
	}
	mmOrderCreate.defaultExpectation.paramPtrs.ctx = &ctx
	mmOrderCreate.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmOrderCreate
}

// ExpectUserIDParam2 sets up expected param userID for LomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) ExpectUserIDParam2(userID int64) *mLomsClientMockOrderCreate {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &LomsClientMockOrderCreateExpectation{}
	}

	if mmOrderCreate.defaultExpectation.params != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Expect")
	}

	if mmOrderCreate.defaultExpectation.paramPtrs == nil {
		mmOrderCreate.defaultExpectation.paramPtrs = &LomsClientMockOrderCreateParamPtrs{}
	}
	mmOrderCreate.defaultExpectation.paramPtrs.userID = &userID
	mmOrderCreate.defaultExpectation.expectationOrigins.originUserID = minimock.CallerInfo(1)

	return mmOrderCreate
}

// ExpectItemsParam3 sets up expected param items for LomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) ExpectItemsParam3(items []model.CartItem) *mLomsClientMockOrderCreate {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &LomsClientMockOrderCreateExpectation{}
	}

	if mmOrderCreate.defaultExpectation.params != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Expect")
	}

	if mmOrderCreate.defaultExpectation.paramPtrs == nil {
		mmOrderCreate.defaultExpectation.paramPtrs = &LomsClientMockOrderCreateParamPtrs{}
	}
	mmOrderCreate.defaultExpectation.paramPtrs.items = &items
	mmOrderCreate.defaultExpectation.expectationOrigins.originItems = minimock.CallerInfo(1)

	return mmOrderCreate
}

// Inspect accepts an inspector function that has same arguments as the LomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) Inspect(f func(ctx context.Context, userID int64, items []model.CartItem)) *mLomsClientMockOrderCreate {
	if mmOrderCreate.mock.inspectFuncOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("Inspect function is already set for LomsClientMock.OrderCreate")
	}

	mmOrderCreate.mock.inspectFuncOrderCreate = f

	return mmOrderCreate
}

// Return sets up results that will be returned by LomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) Return(orderID int64, err error) *LomsClientMock {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &LomsClientMockOrderCreateExpectation{mock: mmOrderCreate.mock}
	}
	mmOrderCreate.defaultExpectation.results = &LomsClientMockOrderCreateResults{orderID, err}
	mmOrderCreate.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmOrderCreate.mock
}

// Set uses given function f to mock the LomsClient.OrderCreate method
func (mmOrderCreate *mLomsClientMockOrderCreate) Set(f func(ctx context.Context, userID int64, items []model.CartItem) (orderID int64, err error)) *LomsClientMock {
	if mmOrderCreate.defaultExpectation != nil {
		mmOrderCreate.mock.t.Fatalf("Default expectation is already set for the LomsClient.OrderCreate method")
	}

	if len(mmOrderCreate.expectations) > 0 {
		mmOrderCreate.mock.t.Fatalf("Some expectations are already set for the LomsClient.OrderCreate method")
	}

	mmOrderCreate.mock.funcOrderCreate = f
	mmOrderCreate.mock.funcOrderCreateOrigin = minimock.CallerInfo(1)
	return mmOrderCreate.mock
}

// When sets expectation for the LomsClient.OrderCreate which will trigger the result defined by the following
// Then helper
func (mmOrderCreate *mLomsClientMockOrderCreate) When(ctx context.Context, userID int64, items []model.CartItem) *LomsClientMockOrderCreateExpectation {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	expectation := &LomsClientMockOrderCreateExpectation{
		mock:               mmOrderCreate.mock,
		params:             &LomsClientMockOrderCreateParams{ctx, userID, items},
		expectationOrigins: LomsClientMockOrderCreateExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmOrderCreate.expectations = append(mmOrderCreate.expectations, expectation)
	return expectation
}

// Then sets up LomsClient.OrderCreate return parameters for the expectation previously defined by the When method
func (e *LomsClientMockOrderCreateExpectation) Then(orderID int64, err error) *LomsClientMock {
	e.results = &LomsClientMockOrderCreateResults{orderID, err}
	return e.mock
}

// Times sets number of times LomsClient.OrderCreate should be invoked
func (mmOrderCreate *mLomsClientMockOrderCreate) Times(n uint64) *mLomsClientMockOrderCreate {
	if n == 0 {
		mmOrderCreate.mock.t.Fatalf("Times of LomsClientMock.OrderCreate mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmOrderCreate.expectedInvocations, n)
	mmOrderCreate.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmOrderCreate
}

func (mmOrderCreate *mLomsClientMockOrderCreate) invocationsDone() bool {
	if len(mmOrderCreate.expectations) == 0 && mmOrderCreate.defaultExpectation == nil && mmOrderCreate.mock.funcOrderCreate == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmOrderCreate.mock.afterOrderCreateCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmOrderCreate.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// OrderCreate implements mm_service.LomsClient
func (mmOrderCreate *LomsClientMock) OrderCreate(ctx context.Context, userID int64, items []model.CartItem) (orderID int64, err error) {
	mm_atomic.AddUint64(&mmOrderCreate.beforeOrderCreateCounter, 1)
	defer mm_atomic.AddUint64(&mmOrderCreate.afterOrderCreateCounter, 1)

	mmOrderCreate.t.Helper()

	if mmOrderCreate.inspectFuncOrderCreate != nil {
		mmOrderCreate.inspectFuncOrderCreate(ctx, userID, items)
	}

	mm_params := LomsClientMockOrderCreateParams{ctx, userID, items}

	// Record call args
	mmOrderCreate.OrderCreateMock.mutex.Lock()
	mmOrderCreate.OrderCreateMock.callArgs = append(mmOrderCreate.OrderCreateMock.callArgs, &mm_params)
	mmOrderCreate.OrderCreateMock.mutex.Unlock()

	for _, e := range mmOrderCreate.OrderCreateMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.orderID, e.results.err
		}
	}

	if mmOrderCreate.OrderCreateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmOrderCreate.OrderCreateMock.defaultExpectation.Counter, 1)
		mm_want := mmOrderCreate.OrderCreateMock.defaultExpectation.params
		mm_want_ptrs := mmOrderCreate.OrderCreateMock.defaultExpectation.paramPtrs

		mm_got := LomsClientMockOrderCreateParams{ctx, userID, items}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmOrderCreate.t.Errorf("LomsClientMock.OrderCreate got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmOrderCreate.OrderCreateMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.userID != nil && !minimock.Equal(*mm_want_ptrs.userID, mm_got.userID) {
				mmOrderCreate.t.Errorf("LomsClientMock.OrderCreate got unexpected parameter userID, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmOrderCreate.OrderCreateMock.defaultExpectation.expectationOrigins.originUserID, *mm_want_ptrs.userID, mm_got.userID, minimock.Diff(*mm_want_ptrs.userID, mm_got.userID))
			}

			if mm_want_ptrs.items != nil && !minimock.Equal(*mm_want_ptrs.items, mm_got.items) {
				mmOrderCreate.t.Errorf("LomsClientMock.OrderCreate got unexpected parameter items, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmOrderCreate.OrderCreateMock.defaultExpectation.expectationOrigins.originItems, *mm_want_ptrs.items, mm_got.items, minimock.Diff(*mm_want_ptrs.items, mm_got.items))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmOrderCreate.t.Errorf("LomsClientMock.OrderCreate got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmOrderCreate.OrderCreateMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmOrderCreate.OrderCreateMock.defaultExpectation.results
		if mm_results == nil {
			mmOrderCreate.t.Fatal("No results are set for the LomsClientMock.OrderCreate")
		}
		return (*mm_results).orderID, (*mm_results).err
	}
	if mmOrderCreate.funcOrderCreate != nil {
		return mmOrderCreate.funcOrderCreate(ctx, userID, items)
	}
	mmOrderCreate.t.Fatalf("Unexpected call to LomsClientMock.OrderCreate. %v %v %v", ctx, userID, items)
	return
}

// OrderCreateAfterCounter returns a count of finished LomsClientMock.OrderCreate invocations
func (mmOrderCreate *LomsClientMock) OrderCreateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmOrderCreate.afterOrderCreateCounter)
}

// OrderCreateBeforeCounter returns a count of LomsClientMock.OrderCreate invocations
func (mmOrderCreate *LomsClientMock) OrderCreateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmOrderCreate.beforeOrderCreateCounter)
}

// Calls returns a list of arguments used in each call to LomsClientMock.OrderCreate.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmOrderCreate *mLomsClientMockOrderCreate) Calls() []*LomsClientMockOrderCreateParams {
	mmOrderCreate.mutex.RLock()

	argCopy := make([]*LomsClientMockOrderCreateParams, len(mmOrderCreate.callArgs))
	copy(argCopy, mmOrderCreate.callArgs)

	mmOrderCreate.mutex.RUnlock()

	return argCopy
}

// MinimockOrderCreateDone returns true if the count of the OrderCreate invocations corresponds
// the number of defined expectations
func (m *LomsClientMock) MinimockOrderCreateDone() bool {
	if m.OrderCreateMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.OrderCreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.OrderCreateMock.invocationsDone()
}

// MinimockOrderCreateInspect logs each unmet expectation
func (m *LomsClientMock) MinimockOrderCreateInspect() {
	for _, e := range m.OrderCreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LomsClientMock.OrderCreate at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterOrderCreateCounter := mm_atomic.LoadUint64(&m.afterOrderCreateCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.OrderCreateMock.defaultExpectation != nil && afterOrderCreateCounter < 1 {
		if m.OrderCreateMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to LomsClientMock.OrderCreate at\n%s", m.OrderCreateMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to LomsClientMock.OrderCreate at\n%s with params: %#v", m.OrderCreateMock.defaultExpectation.expectationOrigins.origin, *m.OrderCreateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcOrderCreate != nil && afterOrderCreateCounter < 1 {
		m.t.Errorf("Expected call to LomsClientMock.OrderCreate at\n%s", m.funcOrderCreateOrigin)
	}

	if !m.OrderCreateMock.invocationsDone() && afterOrderCreateCounter > 0 {
		m.t.Errorf("Expected %d calls to LomsClientMock.OrderCreate at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.OrderCreateMock.expectedInvocations), m.OrderCreateMock.expectedInvocationsOrigin, afterOrderCreateCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *LomsClientMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetStockInspect()

			m.MinimockOrderCreateInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *LomsClientMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *LomsClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetStockDone() &&
		m.MinimockOrderCreateDone()
}
