##############################
# プロバイダ設定
##############################
provider "aws" {
  region = var.region
}

##############################
# セキュリティグループ設定
##############################
resource "aws_security_group" "app_sg" {
  name        = "hackathon-app-sg"
  description = "Allow SSH, HTTP, HTTPS, App"

  # SSH
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTPS
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Go App (8080)
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # アウトバウンド全開放
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "hackathon-app-sg"
  }
}

##############################
# EC2 インスタンス作成
##############################
resource "aws_instance" "app_server" {
  ami           = var.ami
  instance_type = var.instance_type
  key_name      = var.key_name
  
  # セキュリティグループを割り当て
  vpc_security_group_ids = [aws_security_group.app_sg.id]

  # 起動スクリプト読み込み
  user_data = file("user_data.sh")

  tags = {
    Name = "hackathon-app-server"
  }
}

##############################
# Elastic IP (EIP)
##############################
resource "aws_eip" "app_eip" {
  domain = "vpc"
  tags = {
    Name = "hackathon-app-eip"
  }
}

resource "aws_eip_association" "app_assoc" {
  instance_id   = aws_instance.app_server.id
  allocation_id = aws_eip.app_eip.id
}

##############################
# 出力
##############################
output "instance_public_ip" {
  value = aws_eip.app_eip.public_ip
}

output "ssh_command" {
  value = "ssh -i ~/.ssh/${var.key_name} ec2-user@${aws_eip.app_eip.public_ip}"
}