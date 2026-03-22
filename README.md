# ood-sagemaker-adapter

OOD compute adapter for AWS SageMaker Studio. Translates Open OnDemand interactive app requests to SageMaker API calls and returns presigned Studio URLs.

## Commands

| Command | Description |
|---------|-------------|
| `launch` | Create a SageMaker app and return a presigned Studio URL |
| `status <app-id>` | Get OOD-normalized status of a Studio app |
| `delete <app-id>` | Delete a Studio app |

App IDs have the format: `<domain-id>/<user-profile>/<app-type>/<app-name>`

## Usage

```bash
# Launch a Studio session
echo '{"app_name":"session-1","user_name":"jsmith"}' | \
  ood-sagemaker-adapter launch \
    --domain-id d-xxxxxxxx \
    --region us-east-1

# Check status
ood-sagemaker-adapter status d-xxxxxxxx/jsmith/JupyterServer/session-1

# Delete
ood-sagemaker-adapter delete d-xxxxxxxx/jsmith/JupyterServer/session-1
```

## OOD Cluster Config

```yaml
# /etc/ood/config/clusters.d/aws-sagemaker.yml
---
v2:
  metadata:
    title: "AWS SageMaker Studio"
  job:
    adapter: "adapter_script"
    submit_host: "localhost"
    submit:
      script: "/usr/local/lib/ood-adapters/ood-sagemaker-adapter"
      args:
        - launch
        - "--domain-id=d-xxxxxxxx"
        - "--region=us-east-1"
```

## Infrastructure

Terraform in `aws-openondemand` with `adapters_enabled = ["sagemaker"]` provisions:
- SageMaker Domain `ood-<env>`
- Default user profile
- Execution role with `AmazonSageMakerFullAccess`
- IAM policy on the OOD instance role
