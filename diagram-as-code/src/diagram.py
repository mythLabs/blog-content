# diagram.py
from diagrams import Cluster, Diagram, Edge
from diagrams.aws.compute import EC2
from diagrams.aws.database import RDS
from diagrams.aws.network import ELB, CF, ALB, Endpoint
from diagrams.aws.mobile import APIGateway
from diagrams.aws.compute import Lambda, ECS ,EC2
from diagrams.aws.security import WAF, ACM 
from diagrams.aws.general import Client

with Diagram("API Architecture", show=False, direction="LR"):
    Client = Client("user")
    with Cluster("AWS Cloud"):
        cf = CF("AWS Cloudfront")
        acm = ACM("SSL Certificates")
        waf=WAF("WAF")
        with Cluster("API VPC"):
            with Cluster("Public subnet"):
                alb = acm - ALB("Load Balancer")
                alb = waf >> Edge(color="brown") << alb
            with Cluster("Private subnet"):
                endp = Endpoint("VPC endpoint")
        api = APIGateway("Private API")
        with Cluster("Compute VPC"):
            handlers = [Lambda("Lambda"),ECS("ECS"),EC2("EC2")]

        Client >> cf >> alb >> endp >> api >> handlers