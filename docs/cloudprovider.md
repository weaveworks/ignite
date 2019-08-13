# Cloud Provider Instances with KVM support

If you intend to use a cloud provider to test Ignite, you can use the instructions below to provision an instance that satisfies the KVM system requirements described in the [installation guide](./installation.md).

## Amazon Web Services

Amazon EC2 [bare metal instances](https://aws.amazon.com/about-aws/whats-new/2018/05/announcing-general-availability-of-amazon-ec2-bare-metal-instances/) provide direct access to the  Intel® Xeon® Scalable processor and memory resources of the underlying server. These instances are ideal for workloads that require access to the hardware feature set (such as Intel® VT-x), for applications that need to run in non-virtualized environments for licensing or support requirements, or for customers who wish to use their own hypervisor.

Here's a list of instances with KVM support, with pricing (as of July 2019), to help you test Ignite. All the instances listed below are EBS-optimized, with 25 Gigabit available network performance and IPv6 support.

| Family | Type | Pricing (US-West-2) per On Demand Linux Instance Hr | vCPUs | Memory (GiB) | Instance Storage (GB) | 
| ---- | ---- | :----: | :----: | :----: | ---- | 
|Compute optimized | c5.metal | $4.08 | 96 |192 |EBS only | 
| General purpose | m5.metal | $4.608 | 96 | 384 | EBS only |
| General purpose |  m5d.metal | $5.424 | 96 | 384  |4 x 900 (SSD) |
|Memory optimized| r5.metal| $6.048 |96 |768| EBS only| 
|Memory optimized| r5d.metal| $6.912 | 96 |768 |4 x 900 (SSD)| 
|Memory optimized| z1d.metal| $4.464 | 48 |384 |2 x 900 (SSD)|
|Storage optimized| i3.metal| $4.992 | 72 | 512 | 8 x 1900 (SSD) |

Use the AWS console to [launch one of these instances](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/LaunchingAndUsingInstances.html) and [connect to your instance using SSH](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AccessingInstancesLinux.html). Then, follow the instructions in the [installation guide](./installation.md).

## Google Cloud

Source: https://blog.kubernauts.io/ignite-on-google-cloud-5d5228a5ffec

Use Google compute from a custom KVM image so that Ignite can be installed and run easily.

- Login to Google Cloud Console
- Open Google Cloud Shell
- Run the following command to create custom images with KVM enabled

  ```bash
  gcloud compute images create nested-virt \
    --source-image-project=ubuntu-os-cloud \
    --source-image-family=ubuntu-1604-lts \
    --licenses="https://www.googleapis.com/compute/v1/projects/vm-options/global/licenses/enable-vmx"
  ```

- Create an instance with the custom image created

## Microsoft Azure

Azure supports nested virtualizations on Dv3 and Ev3 series virtual machines with 4 or more vCPUs. The smallest sizes supporting Ignite are listed below for development purposes. Larger and more specialized SKUs also support virtualization. There is no special configuration required. Follow the standard Ignite installation instructions for Ubuntu.


| Family | Type | $/hr (West US 2) | $/month | vCPUs | Memory (GiB) | Instance Storage (GB) | 
| ----           | ------ | :----: | :-----: | :-: | :-: | :----: 
|General purpose | D4s_v3 | $0.192 | $140.16 | 4   | 16  | 32GiB  |
|General purpose | D8s_v3 | $0.384 | $280.32 | 8   | 32  | 64GiB  |
|Memory optimized| E4s_v3 | $0.252 | $183.96 | 4   | 32  | 64GiB  | 
|Memory optimized| E8s_v3 | $0.504 | $367.92 | 8   | 64  | 128GiB |


## Packet

As a bare metal provider Packet naturally supports virtualization. For development purposes, using a t1.small.x86 with Ubuntu 18.04 or 19.04 works well.

| Config | $/hr | $/month | CPUs | Memory (GiB) | Instance Storage (GB) |  
| ------ | :----: | :-----: | :-: | :-: | :----: 
| t1.small.x86 | $0.07 | ~$51.1 | 4@2.4Ghz  | 8  | 80GiB   |
| c1.small.x86 | $0.40 | ~$292  | 4@3.5Ghz  | 32 | 120GiB  |


