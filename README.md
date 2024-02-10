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


 # Test hook



 # Argo CD Test
 1. ssm to bastion
    - aws ssm start-session --target i-00de322c50efe6f8b --profile backyard

    https://awstip.com/securely-connect-to-a-private-eks-cluster-using-aws-ssm-session-forwarding-systems-manager-5d0767edea61

    https://medium.com/canisworks/aws-systems-manager-vs-bastion-hosts-for-private-networks-efe9a42f5ad7

    aws ssm start-session --target i-0d7a2ea95ef0ee3ef --document-name AWS-StartPortForwardingSessionToRemoteHost --parameters '{"portNumber":["443"],"localPortNumber":["4443"],"host":["A175A91839F6ED3EA3AA61C9381BD847.gr7.ap-northeast-1.eks.amazonaws.com"]}' --profile backyard