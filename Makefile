port-forward-testkube:
	kubectl -n testkube port-forward svc/testkube-dashboard 8080 &
	kubectl -n testkube port-forward svc/testkube-api-server 8088