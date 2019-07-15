# Usage
This simple application is use for check pod network connectivity.

# Why?
Because on our GKE cluster. We found sometime the node's `kube-proxy` gone. 
And it won't auto recovery. 
It made pod in that node(s) loss network connectivity to anything inside cluster (included DNS and connection to API server).

So that, use this simple application test for network connectivity and can it resolve DNS.
Have another program to check this application readiness status. 
If this pod started over 5 mins. And still not ready. That application will drain that node.
