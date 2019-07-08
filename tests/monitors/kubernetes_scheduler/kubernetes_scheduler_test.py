import pytest

from tests.helpers.metadata import Metadata
from tests.helpers.verify import verify

pytestmark = [pytest.mark.kubernetes]
METADATA = Metadata.from_package("kubernetes/scheduler")


@pytest.skip("only works with real minikube currently")
def test_kubernetes_scheduler(k8s_cluster):
    config = """
        observers:
        - type: k8s-api

        monitors:
        - type: kubernetes-scheduler
          discoveryRule: container_name == "kube-scheduler"
          port: 10251
          extraMetrics: ["*"]
     """
    with k8s_cluster.run_agent(config) as agent:
        verify(agent, METADATA.all_metrics)
