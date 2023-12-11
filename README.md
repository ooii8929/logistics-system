# sre-program-remote-4

# develop env
1. create key pairs in console by manual - "logistics-system"
 chmod 600 logistics-system.pem
 ssh-keygen -R 35.73.82.22
 ssh -v -i logistics-system.pem ubuntu@3.112.101.105
2. create github
 go to AWS console find the deploy IAM user access key and put in github action env


 # CI
 ```
 git add .
 git commit -m "feat / fix / doc : xxxx"
 git push
 ```

 # CD
 ```
 cd infra
 terraform plan -out prod.tfplan
 terraform apply prod.tfplan
 ```

 ```
 Outputs:

eip = ""
repository_url = ""
 ```