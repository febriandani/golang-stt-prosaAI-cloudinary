#mock infra
mockgen -package=mock_infra -source=infra/db.go -destination=infra/mock/db.go

#mock repo
# mockgen -package=mock_repo -source=repo/city.go -destination=repo/mock/city.go