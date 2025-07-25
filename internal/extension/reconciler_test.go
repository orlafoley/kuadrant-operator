//go:build unit

package extension

import (
	"sync"
	"testing"
	"time"

	"github.com/samber/lo"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	machinerycontroller "github.com/kuadrant/policy-machinery/controller"
	"github.com/kuadrant/policy-machinery/machinery"

	v1 "github.com/kuadrant/kuadrant-operator/pkg/extension/grpc/v1"
)

func TestStateAwareDAG(t *testing.T) {
	t.Run("findGateways()", func(t *testing.T) {
		resources := BuildComplexGatewayAPITopology()

		gatewayClasses := lo.Map(resources.GatewayClasses, func(gatewayClass *gwapiv1.GatewayClass, _ int) *machinery.GatewayClass {
			return &machinery.GatewayClass{GatewayClass: gatewayClass}
		})
		gateways := lo.Map(resources.Gateways, func(gateway *gwapiv1.Gateway, _ int) *machinery.Gateway { return &machinery.Gateway{Gateway: gateway} })
		httpRoutes := lo.Map(resources.HTTPRoutes, func(httpRoute *gwapiv1.HTTPRoute, _ int) *machinery.HTTPRoute {
			return &machinery.HTTPRoute{HTTPRoute: httpRoute}
		})
		grpcRoutes := lo.Map(resources.GRPCRoutes, func(grpcRoute *gwapiv1.GRPCRoute, _ int) *machinery.GRPCRoute {
			return &machinery.GRPCRoute{GRPCRoute: grpcRoute}
		})
		tcpRoutes := lo.Map(resources.TCPRoutes, func(tcpRoute *gwapiv1alpha2.TCPRoute, _ int) *machinery.TCPRoute {
			return &machinery.TCPRoute{TCPRoute: tcpRoute}
		})
		tlsRoutes := lo.Map(resources.TLSRoutes, func(tlsRoute *gwapiv1alpha2.TLSRoute, _ int) *machinery.TLSRoute {
			return &machinery.TLSRoute{TLSRoute: tlsRoute}
		})
		udpRoutes := lo.Map(resources.UDPRoutes, func(updRoute *gwapiv1alpha2.UDPRoute, _ int) *machinery.UDPRoute {
			return &machinery.UDPRoute{UDPRoute: updRoute}
		})
		services := lo.Map(resources.Services, func(service *core.Service, _ int) *machinery.Service { return &machinery.Service{Service: service} })

		topology, err := machinery.NewTopology(
			machinery.WithTargetables(gatewayClasses...),
			machinery.WithTargetables(gateways...),
			machinery.WithTargetables(httpRoutes...),
			machinery.WithTargetables(services...),
			machinery.WithTargetables(grpcRoutes...),
			machinery.WithTargetables(tcpRoutes...),
			machinery.WithTargetables(tlsRoutes...),
			machinery.WithTargetables(udpRoutes...),
			machinery.WithLinks(
				machinery.LinkGatewayClassToGatewayFunc(gatewayClasses),
				machinery.LinkGatewayToHTTPRouteFunc(gateways),
				machinery.LinkGatewayToGRPCRouteFunc(gateways),
				machinery.LinkGatewayToTCPRouteFunc(gateways),
				machinery.LinkGatewayToTLSRouteFunc(gateways),
				machinery.LinkGatewayToUDPRouteFunc(gateways),
				machinery.LinkHTTPRouteToServiceFunc(httpRoutes, false),
				machinery.LinkGRPCRouteToServiceFunc(grpcRoutes, false),
				machinery.LinkTCPRouteToServiceFunc(tcpRoutes, false),
				machinery.LinkTLSRouteToServiceFunc(tlsRoutes, false),
				machinery.LinkUDPRouteToServiceFunc(udpRoutes, false),
			),
		)

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		dag := StateAwareDAG{
			topology,
			nil,
		}

		gws, err := dag.FindGatewaysFor([]*v1.TargetRef{{Kind: "Service", Name: "service-1"}})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if len(gws) != 1 {
			t.Fatalf("Expected exactly 1 gateway, got %#v", gws)
		}
		if gws[0].GetMetadata().GetName() != "gateway-1" {
			t.Fatalf("Expected gateway-1, got %s", gws[0].GetMetadata().GetName())
		}

		gws, err = dag.FindGatewaysFor([]*v1.TargetRef{{Kind: "TLSRoute", Name: "tls-route-1"}})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if len(gws) != 2 {
			t.Fatalf("Expected exactly 2 gateways, got %#v", gws)
		}
		if gws[0].GetMetadata().GetName() != "gateway-3" && gws[1].GetMetadata().GetName() != "gateway-3" {
			t.Fatalf("Expected gateway-3")
		}
		if gws[0].GetMetadata().GetName() != "gateway-4" && gws[1].GetMetadata().GetName() != "gateway-4" {
			t.Fatalf("Expected gateway-4")
		}

		gws, err = dag.FindGatewaysFor([]*v1.TargetRef{{Kind: "Service", Name: "service-3"}})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if len(gws) != 3 {
			t.Fatalf("Expected exactly 3 gateways, got %#v", gws)
		}
		if gws[0].GetMetadata().GetName() != "gateway-1" && gws[1].GetMetadata().GetName() != "gateway-1" && gws[2].GetMetadata().GetName() != "gateway-1" {
			t.Fatalf("Expected gateway-1, got %#v", gws[0].GetMetadata().GetName())
		}
		if gws[0].GetMetadata().GetName() != "gateway-2" && gws[1].GetMetadata().GetName() != "gateway-2" && gws[2].GetMetadata().GetName() != "gateway-2" {
			t.Fatalf("Expected gateway-2, got %#v", gws[0].GetMetadata().GetName())
		}
		if gws[0].GetMetadata().GetName() != "gateway-3" && gws[1].GetMetadata().GetName() != "gateway-3" && gws[2].GetMetadata().GetName() != "gateway-3" {
			t.Fatalf("Expected gateway-3, got %#v", gws)
		}
	})

	t.Run("findPolicies()", func(t *testing.T) {
		resources := BuildComplexGatewayAPITopology()

		gatewayClasses := lo.Map(resources.GatewayClasses, func(gatewayClass *gwapiv1.GatewayClass, _ int) *machinery.GatewayClass {
			return &machinery.GatewayClass{GatewayClass: gatewayClass}
		})
		gateways := lo.Map(resources.Gateways, func(gateway *gwapiv1.Gateway, _ int) *machinery.Gateway { return &machinery.Gateway{Gateway: gateway} })
		httpRoutes := lo.Map(resources.HTTPRoutes, func(httpRoute *gwapiv1.HTTPRoute, _ int) *machinery.HTTPRoute {
			return &machinery.HTTPRoute{HTTPRoute: httpRoute}
		})
		grpcRoutes := lo.Map(resources.GRPCRoutes, func(grpcRoute *gwapiv1.GRPCRoute, _ int) *machinery.GRPCRoute {
			return &machinery.GRPCRoute{GRPCRoute: grpcRoute}
		})
		tcpRoutes := lo.Map(resources.TCPRoutes, func(tcpRoute *gwapiv1alpha2.TCPRoute, _ int) *machinery.TCPRoute {
			return &machinery.TCPRoute{TCPRoute: tcpRoute}
		})
		tlsRoutes := lo.Map(resources.TLSRoutes, func(tlsRoute *gwapiv1alpha2.TLSRoute, _ int) *machinery.TLSRoute {
			return &machinery.TLSRoute{TLSRoute: tlsRoute}
		})
		udpRoutes := lo.Map(resources.UDPRoutes, func(updRoute *gwapiv1alpha2.UDPRoute, _ int) *machinery.UDPRoute {
			return &machinery.UDPRoute{UDPRoute: updRoute}
		})
		services := lo.Map(resources.Services, func(service *core.Service, _ int) *machinery.Service { return &machinery.Service{Service: service} })

		policies := []*TestPolicy{
			buildPolicy(func(policy *TestPolicy) {
				policy.Name = "my-policy-1"
				policy.Spec.TargetRef = gwapiv1alpha2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1alpha2.LocalPolicyTargetReference{
						Group: gwapiv1.GroupName,
						Kind:  "Gateway",
						Name:  "gateway-1",
					},
				}
			}),
			buildPolicy(func(policy *TestPolicy) {
				policy.Name = "my-policy-2"
				policy.Spec.TargetRef = gwapiv1alpha2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1alpha2.LocalPolicyTargetReference{
						Group: gwapiv1.GroupName,
						Kind:  "HTTPRoute",
						Name:  "http-route-1",
					},
				}
			}),
		}

		topology, err := machinery.NewTopology(
			machinery.WithTargetables(gatewayClasses...),
			machinery.WithTargetables(gateways...),
			machinery.WithTargetables(httpRoutes...),
			machinery.WithTargetables(services...),
			machinery.WithTargetables(grpcRoutes...),
			machinery.WithTargetables(tcpRoutes...),
			machinery.WithTargetables(tlsRoutes...),
			machinery.WithTargetables(udpRoutes...),
			machinery.WithPolicies(policies...),
			machinery.WithLinks(
				machinery.LinkGatewayClassToGatewayFunc(gatewayClasses),
				machinery.LinkGatewayToHTTPRouteFunc(gateways),
				machinery.LinkGatewayToGRPCRouteFunc(gateways),
				machinery.LinkGatewayToTCPRouteFunc(gateways),
				machinery.LinkGatewayToTLSRouteFunc(gateways),
				machinery.LinkGatewayToUDPRouteFunc(gateways),
				machinery.LinkHTTPRouteToServiceFunc(httpRoutes, false),
				machinery.LinkGRPCRouteToServiceFunc(grpcRoutes, false),
				machinery.LinkTCPRouteToServiceFunc(tcpRoutes, false),
				machinery.LinkTLSRouteToServiceFunc(tlsRoutes, false),
				machinery.LinkUDPRouteToServiceFunc(udpRoutes, false),
			),
		)

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		dag := StateAwareDAG{
			topology,
			nil,
		}

		pols, err := dag.FindPoliciesFor([]*v1.TargetRef{{Kind: "Gateway", Name: "gateway-1"}}, &TestPolicy{})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if len(pols) != 1 {
			t.Fatalf("Expected exactly 1 policy, got %#v", pols)
		}
		if pols[0].GetMetadata().GetName() != "my-policy-1" {
			t.Fatalf("Expected my-policy-1, got %s", pols[0].GetMetadata().GetName())
		}

		pols, err = dag.FindPoliciesFor([]*v1.TargetRef{{Kind: "HTTPRoute", Name: "http-route-1"}}, &TestPolicy{})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if len(pols) != 2 {
			t.Fatalf("Expected exactly 2 policies, got %#v", pols)
		}
		if pols[0].GetMetadata().GetName() != "my-policy-1" && pols[1].GetMetadata().GetName() != "my-policy-1" {
			t.Fatalf("Expected my-policy-1")
		}
		if pols[0].GetMetadata().GetName() != "my-policy-2" && pols[1].GetMetadata().GetName() != "my-policy-2" {
			t.Fatal("Expected my-policy-2")
		}

		pols, err = dag.FindPoliciesFor([]*v1.TargetRef{{Kind: "TestPolicy", Name: "my-policy-2"}}, &TestPolicy{})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if len(pols) != 2 {
			t.Fatalf("Expected exactly 2 policies, got %#v", pols)
		}
		if pols[0].GetMetadata().GetName() != "my-policy-1" && pols[1].GetMetadata().GetName() != "my-policy-1" {
			t.Fatalf("Expected my-policy-1")
		}
		if pols[0].GetMetadata().GetName() != "my-policy-2" && pols[1].GetMetadata().GetName() != "my-policy-2" {
			t.Fatal("Expected my-policy-2")
		}
	})

	t.Run("empty topology", func(t *testing.T) {
		topology, err := machinery.NewTopology()
		if err != nil {
			t.Fatalf("Failed to create empty topology: %v", err)
		}

		dag := StateAwareDAG{
			topology: topology,
			state:    &sync.Map{},
		}

		gateways, err := dag.FindGatewaysFor([]*v1.TargetRef{})
		if err != nil {
			t.Errorf("Expected no error with empty target refs: %v", err)
		}
		if len(gateways) != 0 {
			t.Errorf("Expected empty gateways list, got %d", len(gateways))
		}

		gateways, err = dag.FindGatewaysFor([]*v1.TargetRef{
			{Kind: "Service", Name: "non-existent"},
		})
		if err != nil {
			t.Errorf("Expected no error with non-existent target refs: %v", err)
		}
		if len(gateways) != 0 {
			t.Errorf("Expected empty gateways list for non-existent service, got %d", len(gateways))
		}
	})

	t.Run("policies with various target types", func(t *testing.T) {
		resources := BuildComplexGatewayAPITopology()

		gatewayClasses := lo.Map(resources.GatewayClasses, func(gatewayClass *gwapiv1.GatewayClass, _ int) *machinery.GatewayClass {
			return &machinery.GatewayClass{GatewayClass: gatewayClass}
		})
		gateways := lo.Map(resources.Gateways, func(gateway *gwapiv1.Gateway, _ int) *machinery.Gateway {
			return &machinery.Gateway{Gateway: gateway}
		})
		httpRoutes := lo.Map(resources.HTTPRoutes, func(httpRoute *gwapiv1.HTTPRoute, _ int) *machinery.HTTPRoute {
			return &machinery.HTTPRoute{HTTPRoute: httpRoute}
		})
		services := lo.Map(resources.Services, func(service *core.Service, _ int) *machinery.Service {
			return &machinery.Service{Service: service}
		})

		policies := []*TestPolicy{
			buildPolicy(func(policy *TestPolicy) {
				policy.Name = "gateway-policy"
				policy.Spec.TargetRef = gwapiv1alpha2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1alpha2.LocalPolicyTargetReference{
						Group: gwapiv1.GroupName,
						Kind:  "Gateway",
						Name:  "gateway-1",
					},
				}
			}),
			buildPolicy(func(policy *TestPolicy) {
				policy.Name = "httproute-policy"
				policy.Spec.TargetRef = gwapiv1alpha2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1alpha2.LocalPolicyTargetReference{
						Group: gwapiv1.GroupName,
						Kind:  "HTTPRoute",
						Name:  "http-route-1",
					},
				}
			}),
		}

		topology, err := machinery.NewTopology(
			machinery.WithTargetables(gatewayClasses...),
			machinery.WithTargetables(gateways...),
			machinery.WithTargetables(httpRoutes...),
			machinery.WithTargetables(services...),
			machinery.WithPolicies(policies...),
			machinery.WithLinks(
				machinery.LinkGatewayClassToGatewayFunc(gatewayClasses),
				machinery.LinkGatewayToHTTPRouteFunc(gateways),
				machinery.LinkHTTPRouteToServiceFunc(httpRoutes, false),
			),
		)
		if err != nil {
			t.Fatalf("Failed to create topology: %v", err)
		}

		dag := StateAwareDAG{
			topology: topology,
			state:    &sync.Map{},
		}

		testPolicy := &TestPolicy{}

		foundPolicies, err := dag.FindPoliciesFor([]*v1.TargetRef{
			{Kind: "Gateway", Name: "gateway-1"},
		}, testPolicy)
		if err != nil {
			t.Errorf("Unexpected error finding policies for gateway: %v", err)
		}
		if len(foundPolicies) == 0 {
			t.Error("Expected to find at least one policy for gateway-1")
		}

		foundPolicies, err = dag.FindPoliciesFor([]*v1.TargetRef{
			{Kind: "HTTPRoute", Name: "http-route-1"},
		}, testPolicy)
		if err != nil {
			t.Errorf("Unexpected error finding policies for HTTPRoute: %v", err)
		}

		foundPolicies, err = dag.FindPoliciesFor([]*v1.TargetRef{
			{Kind: "Service", Name: "service-1"},
		}, testPolicy)
		if err != nil {
			t.Errorf("Unexpected error finding policies for Service: %v", err)
		}
	})
}

func TestNilGuardedPointer(t *testing.T) {
	t.Run("set and get", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()

		if ptr.get() != nil {
			t.Errorf("Expected initial value to be nil, got %v", ptr.get())
		}

		value := "test"
		ptr.set(value)

		loaded := ptr.get()
		if loaded == nil {
			t.Error("Expected loaded value to be non-nil")
		} else if *loaded != value {
			t.Errorf("Expected loaded value to be %s, got %s", value, *loaded)
		}
	})

	t.Run("getWait blocks until value is set", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()

		done := make(chan struct{})
		var loaded string

		go func() {
			loaded = ptr.getWait()
			close(done)
		}()

		time.Sleep(100 * time.Millisecond)

		value := "test"
		ptr.set(value)

		select {
		case <-done:
			if loaded != value {
				t.Errorf("Expected loaded value to be %s, got %s", value, loaded)
			}
		case <-time.After(1 * time.Second):
			t.Error("Timed out waiting for getWait to return")
		}
	})

	t.Run("getWait returns immediately if value is already set", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()

		value := "test"
		ptr.set(value)

		start := time.Now()
		loaded := ptr.getWait()
		elapsed := time.Since(start)

		if elapsed > 100*time.Millisecond {
			t.Errorf("Expected getWait to return immediately, took %v", elapsed)
		}

		if loaded != value {
			t.Errorf("Expected loaded value to be %s, got %s", value, loaded)
		}
	})

	t.Run("getWaitWithTimeout returns false on timeout", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()

		start := time.Now()
		_, success := ptr.getWaitWithTimeout(100 * time.Millisecond)
		elapsed := time.Since(start)

		if elapsed < 100*time.Millisecond {
			t.Errorf("Expected getWaitWithTimeout to wait for at least the timeout duration, took %v", elapsed)
		}

		if success {
			t.Error("Expected success to be false on timeout")
		}
	})

	t.Run("getWaitWithTimeout returns true when value is set before timeout", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()

		done := make(chan bool)
		var loaded string

		go func() {
			var success bool
			l, success := ptr.getWaitWithTimeout(1 * time.Second)
			loaded = *l
			done <- success
		}()

		time.Sleep(100 * time.Millisecond)

		value := "test"
		ptr.set(value)

		select {
		case success := <-done:
			if !success {
				t.Error("Expected success to be true when value is set before timeout")
			}
			if loaded != value {
				t.Errorf("Expected loaded value to be %s, got %s", value, loaded)
			}
		case <-time.After(2 * time.Second):
			t.Error("Timed out waiting for getWaitWithTimeout to return")
		}
	})

	t.Run("set sends updates", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()
		channel := ptr.newUpdateChannel()

		if ptr.get() != nil {
			t.Errorf("Expected initial value to be nil, got %v", ptr.get())
		}

		value := "test"
		ptr.set(value)

		loaded := ptr.get()
		if loaded == nil {
			t.Error("Expected loaded value to be non-nil")
		} else if *loaded != value {
			t.Errorf("Expected loaded value to be %s, got %s", value, *loaded)
		}

		go func() {
			ptr.set("updated once")
			ptr.set("updated twice")
		}()

		one, two := <-channel, <-channel
		if one != "updated once" {
			t.Errorf("Expected update to be `updated once`, got `%s`", one)
		}
		if two != "updated twice" {
			t.Errorf("Expected update to be `updated twice`, got `%s`", two)
		}
	})

	t.Run("BlockingDAG variable", func(t *testing.T) {
		if BlockingDAG.get() != nil {
			t.Error("Expected initial BlockingDAG to be nil")
		}

		dag := StateAwareDAG{
			topology: nil,
			state:    &sync.Map{},
		}

		BlockingDAG.set(dag)

		loaded := BlockingDAG.get()
		if loaded == nil {
			t.Error("Expected loaded BlockingDAG to be non-nil")
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()

		var wg sync.WaitGroup
		results := make([]string, 10)

		for i := range 10 {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				results[index] = ptr.getWait()
			}(i)
		}

		go func() {
			time.Sleep(10 * time.Millisecond)
			ptr.set("test-value")
		}()

		wg.Wait()

		for i, result := range results {
			if result != "test-value" {
				t.Errorf("Goroutine %d got '%s', expected 'test-value'", i, result)
			}
		}
	})

	t.Run("getWaitWithTimeout success", func(t *testing.T) {
		ptr := newNilGuardedPointer[int]()

		go func() {
			time.Sleep(10 * time.Millisecond)
			ptr.set(42)
		}()

		value, ok := ptr.getWaitWithTimeout(100 * time.Millisecond)
		if !ok {
			t.Error("Expected timeout to succeed")
		}
		if value == nil || *value != 42 {
			t.Errorf("Expected value 42, got %v", value)
		}
	})

	t.Run("getWaitWithTimeout timeout", func(t *testing.T) {
		ptr := newNilGuardedPointer[int]()

		start := time.Now()
		value, ok := ptr.getWaitWithTimeout(50 * time.Millisecond)
		duration := time.Since(start)

		if ok {
			t.Error("Expected timeout to fail")
		}
		if value != nil {
			t.Errorf("Expected nil value on timeout, got %v", value)
		}
		if duration < 40*time.Millisecond || duration > 100*time.Millisecond {
			t.Errorf("Expected timeout around 50ms, got %v", duration)
		}
	})

	t.Run("multiple sets and updates", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()

		ch1 := ptr.newUpdateChannel()
		ch2 := ptr.newUpdateChannel()

		ptr.set("first")

		var wg sync.WaitGroup
		results := make([]string, 2)

		wg.Add(2)
		go func() {
			defer wg.Done()
			select {
			case val := <-ch1:
				results[0] = val
			case <-time.After(200 * time.Millisecond):
				results[0] = "timeout"
			}
		}()

		go func() {
			defer wg.Done()
			select {
			case val := <-ch2:
				results[1] = val
			case <-time.After(200 * time.Millisecond):
				results[1] = "timeout"
			}
		}()

		ptr.set("second")

		wg.Wait()

		if results[0] != "second" {
			t.Errorf("Expected 'second' on update channel 1, got '%s'", results[0])
		}
		if results[1] != "second" {
			t.Errorf("Expected 'second' on update channel 2, got '%s'", results[1])
		}
	})

	t.Run("get without waiting", func(t *testing.T) {
		ptr := newNilGuardedPointer[string]()

		if val := ptr.get(); val != nil {
			t.Errorf("Expected nil initially, got %v", val)
		}

		ptr.set("test")

		if val := ptr.get(); val == nil || *val != "test" {
			t.Errorf("Expected 'test', got %v", val)
		}
	})
}

func TestConversionFunctions(t *testing.T) {
	t.Run("toGw conversion", func(t *testing.T) {
		gateway := machinery.Gateway{
			Gateway: &gwapiv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-gateway",
					Namespace: "test-namespace",
				},
				Spec: gwapiv1.GatewaySpec{
					GatewayClassName: "test-class",
					Listeners: []gwapiv1.Listener{
						{
							Name:     "http",
							Hostname: ptr.To[gwapiv1.Hostname]("example.com"),
							Port:     80,
							Protocol: gwapiv1.HTTPProtocolType,
						},
						{
							Name:     "https",
							Port:     443,
							Protocol: gwapiv1.HTTPSProtocolType,
							// No hostname to test nil handling
						},
					},
				},
			},
		}

		result := toGw(gateway)

		if result.Metadata.Name != "test-gateway" {
			t.Errorf("Expected name 'test-gateway', got %s", result.Metadata.Name)
		}
		if result.Metadata.Namespace != "test-namespace" {
			t.Errorf("Expected namespace 'test-namespace', got %s", result.Metadata.Namespace)
		}
		if result.Spec.GatewayClassName != "test-class" {
			t.Errorf("Expected gateway class 'test-class', got %s", result.Spec.GatewayClassName)
		}
		if len(result.Spec.Listeners) != 2 {
			t.Errorf("Expected 2 listeners, got %d", len(result.Spec.Listeners))
		}
		if result.Spec.Listeners[0].Hostname != "example.com" {
			t.Errorf("Expected hostname 'example.com', got %s", result.Spec.Listeners[0].Hostname)
		}
		if result.Spec.Listeners[1].Hostname != "" {
			t.Errorf("Expected empty hostname for second listener, got %s", result.Spec.Listeners[1].Hostname)
		}
	})

	t.Run("toPolicy conversion", func(t *testing.T) {
		policy := &TestPolicy{
			TypeMeta: metav1.TypeMeta{
				Kind:       "TestPolicy",
				APIVersion: "test.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-policy",
				Namespace: "test-namespace",
			},
			Spec: TestPolicySpec{
				TargetRef: gwapiv1alpha2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1alpha2.LocalPolicyTargetReference{
						Group: gwapiv1.GroupName,
						Kind:  "Gateway",
						Name:  "test-gateway",
					},
				},
			},
		}

		result := toPolicy(policy)

		if result.Metadata.Name != "test-policy" {
			t.Errorf("Expected name 'test-policy', got %s", result.Metadata.Name)
		}
		if result.Metadata.Namespace != "test-namespace" {
			t.Errorf("Expected namespace 'test-namespace', got %s", result.Metadata.Namespace)
		}
		if result.Metadata.Kind != "TestPolicy" {
			t.Errorf("Expected kind 'TestPolicy', got %s", result.Metadata.Kind)
		}
		if len(result.TargetRefs) != 1 {
			t.Errorf("Expected 1 target ref, got %d", len(result.TargetRefs))
		}
		if result.TargetRefs[0].Name != "test-gateway" {
			t.Errorf("Expected target ref name 'test-gateway', got %s", result.TargetRefs[0].Name)
		}
	})

	t.Run("toTargetRefs conversion", func(t *testing.T) {
		policy := &TestPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-target",
				Namespace: "test-namespace",
			},
			Spec: TestPolicySpec{
				TargetRef: gwapiv1alpha2.LocalPolicyTargetReferenceWithSectionName{
					LocalPolicyTargetReference: gwapiv1alpha2.LocalPolicyTargetReference{
						Group: gwapiv1.GroupName,
						Kind:  "Gateway",
						Name:  "test-gateway",
					},
				},
			},
		}

		targetRefs := policy.GetTargetRefs()
		result := toTargetRefs(targetRefs)

		if len(result) != 1 {
			t.Errorf("Expected 1 target ref, got %d", len(result))
		}
		if result[0].Name != "test-gateway" {
			t.Errorf("Expected name 'test-gateway', got %s", result[0].Name)
		}
		if result[0].Kind != "Gateway" {
			t.Errorf("Expected kind 'Gateway', got %s", result[0].Kind)
		}
		if result[0].Group != gwapiv1.GroupName {
			t.Errorf("Expected group '%s', got %s", gwapiv1.GroupName, result[0].Group)
		}
	})

	t.Run("toListeners conversion", func(t *testing.T) {
		listeners := []gwapiv1.Listener{
			{
				Name:     "http",
				Hostname: ptr.To[gwapiv1.Hostname]("example.com"),
				Port:     80,
				Protocol: gwapiv1.HTTPProtocolType,
			},
			{
				Name:     "https",
				Port:     443,
				Protocol: gwapiv1.HTTPSProtocolType,
				// No hostname to test nil handling
			},
		}

		result := toListeners(listeners)

		if len(result) != 2 {
			t.Errorf("Expected 2 listeners, got %d", len(result))
		}
		if result[0].Hostname != "example.com" {
			t.Errorf("Expected hostname 'example.com', got %s", result[0].Hostname)
		}
		if result[1].Hostname != "" {
			t.Errorf("Expected empty hostname for second listener, got %s", result[1].Hostname)
		}
	})
}

func BuildGatewayClass(f ...func(*gwapiv1.GatewayClass)) *gwapiv1.GatewayClass {
	gc := &gwapiv1.GatewayClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: gwapiv1.GroupVersion.String(),
			Kind:       "GatewayClass",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-gateway-class",
		},
		Spec: gwapiv1.GatewayClassSpec{
			ControllerName: gwapiv1.GatewayController("my-gateway-controller"),
		},
	}
	for _, fn := range f {
		fn(gc)
	}
	return gc
}

func BuildGateway(f ...func(*gwapiv1.Gateway)) *gwapiv1.Gateway {
	g := &gwapiv1.Gateway{
		TypeMeta: metav1.TypeMeta{
			APIVersion: gwapiv1.GroupVersion.String(),
			Kind:       "Gateway",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-gateway",
			Namespace: "my-namespace",
		},
		Spec: gwapiv1.GatewaySpec{
			GatewayClassName: "my-gateway-class",
			Listeners: []gwapiv1.Listener{
				{
					Name:     "my-listener",
					Port:     80,
					Protocol: "HTTP",
				},
			},
		},
	}
	for _, fn := range f {
		fn(g)
	}
	return g
}

func BuildHTTPRoute(f ...func(*gwapiv1.HTTPRoute)) *gwapiv1.HTTPRoute {
	r := &gwapiv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: gwapiv1.GroupVersion.String(),
			Kind:       "HTTPRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-http-route",
			Namespace: "my-namespace",
		},
		Spec: gwapiv1.HTTPRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{
					{
						Name: "my-gateway",
					},
				},
			},
			Rules: []gwapiv1.HTTPRouteRule{
				{
					BackendRefs: []gwapiv1.HTTPBackendRef{BuildHTTPBackendRef()},
				},
			},
		},
	}
	for _, fn := range f {
		fn(r)
	}
	return r
}

func BuildHTTPBackendRef(f ...func(*gwapiv1.BackendObjectReference)) gwapiv1.HTTPBackendRef {
	return gwapiv1.HTTPBackendRef{
		BackendRef: BuildBackendRef(f...),
	}
}

func BuildService(f ...func(*core.Service)) *core.Service {
	s := &core.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: core.SchemeGroupVersion.String(),
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-service",
			Namespace: "my-namespace",
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
			Selector: map[string]string{
				"app": "my-app",
			},
		},
	}
	for _, fn := range f {
		fn(s)
	}
	return s
}

func BuildGRPCRoute(f ...func(*gwapiv1.GRPCRoute)) *gwapiv1.GRPCRoute {
	r := &gwapiv1.GRPCRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: gwapiv1.GroupVersion.String(),
			Kind:       "GRPCRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-grpc-route",
			Namespace: "my-namespace",
		},
		Spec: gwapiv1.GRPCRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{
					{
						Name: "my-gateway",
					},
				},
			},
			Rules: []gwapiv1.GRPCRouteRule{
				{
					BackendRefs: []gwapiv1.GRPCBackendRef{BuildGRPCBackendRef()},
				},
			},
		},
	}
	for _, fn := range f {
		fn(r)
	}

	return r
}

func BuildGRPCBackendRef(f ...func(*gwapiv1.BackendObjectReference)) gwapiv1.GRPCBackendRef {
	return gwapiv1.GRPCBackendRef{
		BackendRef: BuildBackendRef(f...),
	}
}

func BuildBackendRef(f ...func(*gwapiv1.BackendObjectReference)) gwapiv1.BackendRef {
	bor := &gwapiv1.BackendObjectReference{
		Name: "my-service",
	}
	for _, fn := range f {
		fn(bor)
	}
	return gwapiv1.BackendRef{
		BackendObjectReference: *bor,
	}
}

func BuildTCPRoute(f ...func(route *gwapiv1alpha2.TCPRoute)) *gwapiv1alpha2.TCPRoute {
	r := &gwapiv1alpha2.TCPRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: gwapiv1alpha2.GroupVersion.String(),
			Kind:       "TCPRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-tcp-route",
			Namespace: "my-namespace",
		},
		Spec: gwapiv1alpha2.TCPRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{
					{
						Name: "my-gateway",
					},
				},
			},
			Rules: []gwapiv1alpha2.TCPRouteRule{
				{
					BackendRefs: []gwapiv1.BackendRef{BuildBackendRef()},
				},
			},
		},
	}
	for _, fn := range f {
		fn(r)
	}

	return r
}

func BuildTLSRoute(f ...func(route *gwapiv1alpha2.TLSRoute)) *gwapiv1alpha2.TLSRoute {
	r := &gwapiv1alpha2.TLSRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: gwapiv1alpha2.GroupVersion.String(),
			Kind:       "TLSRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-tls-route",
			Namespace: "my-namespace",
		},
		Spec: gwapiv1alpha2.TLSRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{
					{
						Name: "my-gateway",
					},
				},
			},
			Rules: []gwapiv1alpha2.TLSRouteRule{
				{
					BackendRefs: []gwapiv1.BackendRef{BuildBackendRef()},
				},
			},
		},
	}
	for _, fn := range f {
		fn(r)
	}

	return r
}

func BuildUDPRoute(f ...func(route *gwapiv1alpha2.UDPRoute)) *gwapiv1alpha2.UDPRoute {
	r := &gwapiv1alpha2.UDPRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: gwapiv1alpha2.GroupVersion.String(),
			Kind:       "UDPRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-udp-route",
			Namespace: "my-namespace",
		},
		Spec: gwapiv1alpha2.UDPRouteSpec{
			CommonRouteSpec: gwapiv1.CommonRouteSpec{
				ParentRefs: []gwapiv1.ParentReference{
					{
						Name: "my-gateway",
					},
				},
			},
			Rules: []gwapiv1alpha2.UDPRouteRule{
				{
					BackendRefs: []gwapiv1.BackendRef{BuildBackendRef()},
				},
			},
		},
	}
	for _, fn := range f {
		fn(r)
	}

	return r
}

type GatewayAPIResources struct {
	GatewayClasses []*gwapiv1.GatewayClass
	Gateways       []*gwapiv1.Gateway
	HTTPRoutes     []*gwapiv1.HTTPRoute
	GRPCRoutes     []*gwapiv1.GRPCRoute
	TCPRoutes      []*gwapiv1alpha2.TCPRoute
	TLSRoutes      []*gwapiv1alpha2.TLSRoute
	UDPRoutes      []*gwapiv1alpha2.UDPRoute
	Services       []*core.Service
}

// BuildComplexGatewayAPITopology returns a set of Gateway API resources organized :
//
//	                                          ┌────────────────┐                                                                        ┌────────────────┐
//	                                          │ gatewayclass-1 │                                                                        │ gatewayclass-2 │
//	                                          └────────────────┘                                                                        └────────────────┘
//	                                                  ▲                                                                                         ▲
//	                                                  │                                                                                         │
//	                        ┌─────────────────────────┼──────────────────────────┐                                                 ┌────────────┴─────────────┐
//	                        │                         │                          │                                                 │                          │
//	        ┌───────────────┴───────────────┐ ┌───────┴────────┐ ┌───────────────┴───────────────┐                  ┌──────────────┴────────────────┐ ┌───────┴────────┐
//	        │           gateway-1           │ │   gateway-2    │ │           gateway-3           │                  │           gateway-4           │ │   gateway-5    │
//	        │                               │ │                │ │                               │                  │                               │ │                │
//	        │ ┌────────────┐ ┌────────────┐ │ │ ┌────────────┐ │ │ ┌────────────┐ ┌────────────┐ │                  │ ┌────────────┐ ┌────────────┐ │ │ ┌────────────┐ │
//	        │ │ listener-1 │ │ listener-2 │ │ │ │ listener-1 │ │ │ │ listener-1 │ │ listener-2 │ │                  │ │ listener-1 │ │ listener-2 │ │ │ │ listener-1 │ │
//	        │ └────────────┘ └────────────┘ │ │ └────────────┘ │ │ └────────────┘ └────────────┘ │                  │ └────────────┘ └────────────┘ │ │ └────────────┘ │
//	        │                        ▲      │ │      ▲         │ │                               │                  │                               │ │                │
//	        └────────────────────────┬──────┘ └──────┬─────────┘ └───────────────────────────────┘                  └───────────────────────────────┘ └────────────────┘
//	                    ▲            │               │       ▲                    ▲            ▲                            ▲           ▲                        ▲
//	                    │            │               │       │                    │            │                            │           │                        │
//	                    │            └───────┬───────┘       │                    │            └──────────────┬─────────────┘           │                        │
//	                    │                    │               │                    │                           │                         │                        │
//	        ┌───────────┴───────────┐ ┌──────┴───────┐ ┌─────┴────────┐ ┌─────────┴─────────────┐ ┌───────────┴───────────┐ ┌───────────┴───────────┐      ┌─────┴────────┐
//	        │     http-route-1      │ │ http-route-2 │ │ http-route-3 │ │     udp-route-1       │ │      tls-route-1      │ │     tcp-route-1       │      │ grpc-route-1 │
//	        │                       │ │              │ │              │ │                       │ │                       │ │                       │      │              │
//	        │ ┌────────┐ ┌────────┐ │ │ ┌────────┐   │ │  ┌────────┐  │ │ ┌────────┐ ┌────────┐ │ │ ┌────────┐ ┌────────┐ │ │ ┌────────┐ ┌────────┐ │      │ ┌────────┐   │
//	        │ │ rule-1 │ │ rule-2 │ │ │ │ rule-1 │   │ │  │ rule-1 │  │ │ │ rule-1 │ │ rule-2 │ │ │ │ rule-1 │ │ rule-2 │ │ │ │ rule-1 │ │ rule-2 │ │      │ │ rule-1 │   │
//	        │ └────┬───┘ └─────┬──┘ │ │ └────┬───┘   │ │  └───┬────┘  │ │ └─┬──────┘ └───┬────┘ │ │ └───┬────┘ └────┬───┘ │ │ └─┬────┬─┘ └────┬───┘ │      │ └────┬───┘   │
//	        │      │           │    │ │      │       │ │      │       │ │   │            │      │ │     │           │     │ │   │    │        │     │      │      │       │
//	        └──────┼───────────┼────┘ └──────┼───────┘ └──────┼───────┘ └───┼────────────┼──────┘ └─────┼───────────┼─────┘ └───┼────┼────────┼─────┘      └──────┼───────┘
//	               │           │             │                │             │            │              │           │           │    │        │                   │
//	               │           │             └────────────────┤             │            │              └───────────┴───────────┘    │        │                   │
//	               ▼           ▼                              │             │            │                          ▼                ▼        │                   ▼
//	┌───────────────────────┐ ┌────────────┐          ┌───────┴─────────────┴───┐  ┌─────┴──────┐             ┌────────────┐        ┌─────────┴──┐          ┌────────────┐
//	│                       │ │            │          │       ▼             ▼   │  │     ▼      │             │            │        │         ▼  │          │            │
//	│ ┌────────┐ ┌────────┐ │ │ ┌────────┐ │          │   ┌────────┐ ┌────────┐ │  │ ┌────────┐ │             │ ┌────────┐ │        │ ┌────────┐ │          │ ┌────────┐ │
//	│ │ port-1 │ │ port-2 │ │ │ │ port-1 │ │          │   │ port-1 │ │ port-2 │ │  │ │ port-1 │ │             │ │ port-1 │ │        │ │ port-1 │ │          │ │ port-1 │ │
//	│ └────────┘ └────────┘ │ │ └────────┘ │          │   └────────┘ └────────┘ │  │ └────────┘ │             │ └────────┘ │        │ └────────┘ │          │ └────────┘ │
//	│                       │ │            │          │                         │  │            │             │            │        │            │          │            │
//	│       service-1       │ │  service-2 │          │         service-3       │  │  service-4 │             │  service-5 │        │  service-6 │          │  service-7 │
//	└───────────────────────┘ └────────────┘          └─────────────────────────┘  └────────────┘             └────────────┘        └────────────┘          └────────────┘
func BuildComplexGatewayAPITopology(funcs ...func(*GatewayAPIResources)) GatewayAPIResources {
	t := GatewayAPIResources{
		GatewayClasses: []*gwapiv1.GatewayClass{
			BuildGatewayClass(func(gc *gwapiv1.GatewayClass) { gc.Name = "gatewayclass-1" }),
			BuildGatewayClass(func(gc *gwapiv1.GatewayClass) { gc.Name = "gatewayclass-2" }),
		},
		Gateways: []*gwapiv1.Gateway{
			BuildGateway(func(g *gwapiv1.Gateway) {
				g.Name = "gateway-1"
				g.Spec.GatewayClassName = "gatewayclass-1"
				g.Spec.Listeners[0].Name = "listener-1"
				g.Spec.Listeners = append(g.Spec.Listeners, gwapiv1.Listener{
					Name:     "listener-2",
					Port:     443,
					Protocol: "HTTPS",
				})
			}),
			BuildGateway(func(g *gwapiv1.Gateway) {
				g.Name = "gateway-2"
				g.Spec.GatewayClassName = "gatewayclass-1"
				g.Spec.Listeners[0].Name = "listener-1"
			}),
			BuildGateway(func(g *gwapiv1.Gateway) {
				g.Name = "gateway-3"
				g.Spec.GatewayClassName = "gatewayclass-1"
				g.Spec.Listeners[0].Name = "listener-1"
				g.Spec.Listeners = append(g.Spec.Listeners, gwapiv1.Listener{
					Name:     "listener-2",
					Port:     443,
					Protocol: "HTTPS",
				})
			}),
			BuildGateway(func(g *gwapiv1.Gateway) {
				g.Name = "gateway-4"
				g.Spec.GatewayClassName = "gatewayclass-2"
				g.Spec.Listeners[0].Name = "listener-1"
				g.Spec.Listeners = append(g.Spec.Listeners, gwapiv1.Listener{
					Name:     "listener-2",
					Port:     443,
					Protocol: "HTTPS",
				})
			}),
			BuildGateway(func(g *gwapiv1.Gateway) {
				g.Name = "gateway-5"
				g.Spec.GatewayClassName = "gatewayclass-2"
				g.Spec.Listeners[0].Name = "listener-1"
			}),
		},
		HTTPRoutes: []*gwapiv1.HTTPRoute{
			BuildHTTPRoute(func(r *gwapiv1.HTTPRoute) {
				r.Name = "http-route-1"
				r.Spec.ParentRefs[0].Name = "gateway-1"
				r.Spec.Rules = []gwapiv1.HTTPRouteRule{
					{ // rule-1
						BackendRefs: []gwapiv1.HTTPBackendRef{BuildHTTPBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
							backendRef.Name = "service-1"
						})},
					},
					{ // rule-2
						BackendRefs: []gwapiv1.HTTPBackendRef{BuildHTTPBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
							backendRef.Name = "service-2"
						})},
					},
				}
			}),
			BuildHTTPRoute(func(r *gwapiv1.HTTPRoute) {
				r.Name = "http-route-2"
				r.Spec.ParentRefs = []gwapiv1.ParentReference{
					{
						Name:        "gateway-1",
						SectionName: ptr.To(gwapiv1.SectionName("listener-2")),
					},
					{
						Name:        "gateway-2",
						SectionName: ptr.To(gwapiv1.SectionName("listener-1")),
					},
				}
				r.Spec.Rules[0].BackendRefs[0] = BuildHTTPBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
					backendRef.Name = "service-3"
					backendRef.Port = ptr.To(gwapiv1.PortNumber(80)) // port-1
				})
			}),
			BuildHTTPRoute(func(r *gwapiv1.HTTPRoute) {
				r.Name = "http-route-3"
				r.Spec.ParentRefs[0].Name = "gateway-2"
				r.Spec.Rules[0].BackendRefs[0] = BuildHTTPBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
					backendRef.Name = "service-3"
					backendRef.Port = ptr.To(gwapiv1.PortNumber(80)) // port-1
				})
			}),
		},
		Services: []*core.Service{
			BuildService(func(s *core.Service) {
				s.Name = "service-1"
				s.Spec.Ports[0].Name = "port-1"
				s.Spec.Ports = append(s.Spec.Ports, core.ServicePort{
					Name: "port-2",
					Port: 443,
				})
			}),
			BuildService(func(s *core.Service) {
				s.Name = "service-2"
				s.Spec.Ports[0].Name = "port-1"
			}),
			BuildService(func(s *core.Service) {
				s.Name = "service-3"
				s.Spec.Ports[0].Name = "port-1"
				s.Spec.Ports = append(s.Spec.Ports, core.ServicePort{
					Name: "port-2",
					Port: 443,
				})
			}),
			BuildService(func(s *core.Service) {
				s.Name = "service-4"
				s.Spec.Ports[0].Name = "port-1"
			}),
			BuildService(func(s *core.Service) {
				s.Name = "service-5"
				s.Spec.Ports[0].Name = "port-1"
			}),
			BuildService(func(s *core.Service) {
				s.Name = "service-6"
				s.Spec.Ports[0].Name = "port-1"
			}),
			BuildService(func(s *core.Service) {
				s.Name = "service-7"
				s.Spec.Ports[0].Name = "port-1"
			}),
		},
		GRPCRoutes: []*gwapiv1.GRPCRoute{
			BuildGRPCRoute(func(r *gwapiv1.GRPCRoute) {
				r.Name = "grpc-route-1"
				r.Spec.ParentRefs[0].Name = "gateway-5"
				r.Spec.Rules[0].BackendRefs[0] = BuildGRPCBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
					backendRef.Name = "service-7"
				})
			}),
		},
		TCPRoutes: []*gwapiv1alpha2.TCPRoute{
			BuildTCPRoute(func(r *gwapiv1alpha2.TCPRoute) {
				r.Name = "tcp-route-1"
				r.Spec.ParentRefs[0].Name = "gateway-4"
				r.Spec.Rules = []gwapiv1alpha2.TCPRouteRule{
					{ // rule-1
						BackendRefs: []gwapiv1.BackendRef{
							BuildBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
								backendRef.Name = "service-5"
							}),
							BuildBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
								backendRef.Name = "service-6"
							}),
						},
					},
					{ // rule-2
						BackendRefs: []gwapiv1.BackendRef{BuildBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
							backendRef.Name = "service-6"
							backendRef.Port = ptr.To(gwapiv1.PortNumber(80)) // port-1
						})},
					},
				}
			}),
		},
		TLSRoutes: []*gwapiv1alpha2.TLSRoute{
			BuildTLSRoute(func(r *gwapiv1alpha2.TLSRoute) {
				r.Name = "tls-route-1"
				r.Spec.ParentRefs[0].Name = "gateway-3"
				r.Spec.ParentRefs = append(r.Spec.ParentRefs, gwapiv1.ParentReference{Name: "gateway-4"})
				r.Spec.Rules = []gwapiv1alpha2.TLSRouteRule{
					{ // rule-1
						BackendRefs: []gwapiv1.BackendRef{BuildBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
							backendRef.Name = "service-5"
						})},
					},
					{ // rule-2
						BackendRefs: []gwapiv1.BackendRef{BuildBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
							backendRef.Name = "service-5"
						})},
					},
				}
			}),
		},
		UDPRoutes: []*gwapiv1alpha2.UDPRoute{
			BuildUDPRoute(func(r *gwapiv1alpha2.UDPRoute) {
				r.Name = "udp-route-1"
				r.Spec.ParentRefs[0].Name = "gateway-3"
				r.Spec.Rules = []gwapiv1alpha2.UDPRouteRule{
					{ // rule-1
						BackendRefs: []gwapiv1.BackendRef{BuildBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
							backendRef.Name = "service-3"
							backendRef.Port = ptr.To(gwapiv1.PortNumber(443)) // port-2
						})},
					},
					{ // rule-2
						BackendRefs: []gwapiv1.BackendRef{BuildBackendRef(func(backendRef *gwapiv1.BackendObjectReference) {
							backendRef.Name = "service-4"
							backendRef.Port = ptr.To(gwapiv1.PortNumber(80)) // port-1
						})},
					},
				}
			}),
		},
	}
	for _, f := range funcs {
		f(&t)
	}
	return t
}

type TestPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TestPolicySpec `json:"spec"`
}

func (p *TestPolicy) DeepCopyObject() runtime.Object {
	return nil
}

type TestPolicySpec struct {
	TargetRef gwapiv1alpha2.LocalPolicyTargetReferenceWithSectionName `json:"targetRef"`
}

var _ machinery.Policy = &TestPolicy{}
var _ machinerycontroller.Object = &TestPolicy{}

func (p *TestPolicy) GetLocator() string {
	return machinery.LocatorFromObject(p)
}

func (p *TestPolicy) GetTargetRefs() []machinery.PolicyTargetReference {
	return []machinery.PolicyTargetReference{
		machinery.LocalPolicyTargetReferenceWithSectionName{
			LocalPolicyTargetReferenceWithSectionName: p.Spec.TargetRef,
			PolicyNamespace: p.Namespace,
		},
	}
}

func (p *TestPolicy) GetMergeStrategy() machinery.MergeStrategy {
	return machinery.DefaultMergeStrategy
}

func (p *TestPolicy) Merge(_ machinery.Policy) machinery.Policy {
	return &TestPolicy{
		Spec: p.Spec,
	}
}

func buildPolicy(f ...func(*TestPolicy)) *TestPolicy {
	p := &TestPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "test/v1",
			Kind:       "TestPolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-policy",
			Namespace: "my-namespace",
		},
		Spec: TestPolicySpec{
			TargetRef: gwapiv1alpha2.LocalPolicyTargetReferenceWithSectionName{
				LocalPolicyTargetReference: gwapiv1alpha2.LocalPolicyTargetReference{
					Group: gwapiv1.Group(core.SchemeGroupVersion.Group),
					Kind:  "Service",
					Name:  "my-service",
				},
			},
		},
	}
	for _, fn := range f {
		fn(p)
	}
	return p
}
