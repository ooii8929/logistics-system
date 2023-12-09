data "aws_ami" "ubuntu_2204" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners     = ["099720109477"] # Canonical (the company behind Ubuntu)
}


resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu_2204.id
  subnet_id   = aws_subnet.main.id
  associate_public_ip_address = true
  security_groups = [aws_security_group.allow_tls.id]
  key_name = "logistics-system"
  instance_type = "t3.micro"
  user_data = file("./init-instance.sh")
}

resource "aws_security_group" "allow_tls" {
  name        = "allow_tls"
  description = "Allow TLS inbound traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

}

output "instance_public_ip" {
  value       = aws_instance.web.public_ip
  description = "The public IP address of the web instance"
}