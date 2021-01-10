# vorto devops Challenge

Directory Structure

- ./grafana-dashboard --- custom grafana dashboard based on application metrics

- ./helm -- application helm chart

- ./monitoring-resources --- prometheus operator deployment resources to monitor application and kubernetes cluster

- ./postgresql -- postgresql helm chart

- ./terra-gke -- google kubernetes engine terraform automtion files

- .gitlab-ci.yml -- gitlab application automation file

- main.go -- application returns invalid deliverys

- Dockerfile -- multistage application Dockerfile

#Notes

- Application has 2 endpoints. 1 (/) for invalid deliverys, 1 for (/metrics) application related prometheus metrics

- within custom resource definition (service monitor) application metric endpoint automatically detects and monitored by prometheus

