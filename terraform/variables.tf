variable "region" {
  description = "AWS Region"
  type        = string
  default     = "ap-northeast-1"
}

variable "ami" {
  description = "Amazon Linux 2023 AMI (ap-northeast-1)"
  type        = string
  default     = "ami-03852a41f1e05c8e4" 
}

variable "instance_type" {
  description = "Instance Type"
  type        = string
  default     = "t3.medium"
}

variable "key_name" {
  description = "既存のキーペア名を指定してください (例: hackathon_deploy_key)"
  type        = string
  default = "meijo-hackathon-2025"
}