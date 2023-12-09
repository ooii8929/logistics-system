# sre-program-remote-4

1. create key pairs in console by manual - "logistics-system"
 chmod 600 logistics-system.pem
 ssh -v -i logistics-system.pem ubuntu@3.112.101.105


 go to AWS console find the deploy IAM user access key and put in github action env