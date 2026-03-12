# Homelab Dashboard

통합 홈랩 관리 대시보드 - K8s 클러스터 및 서비스 모니터링/관리 UI

## 기술 스택
- **Backend**: Go 1.22+, Echo framework
- **Frontend**: Vue 3, Vite, TypeScript, Tailwind CSS
- **API**: REST (추후 WebSocket for realtime)

## 프로젝트 구조
```
homelab-dashboard/
├── backend/           # Go API 서버
│   ├── cmd/           # 엔트리포인트
│   ├── internal/      # 내부 패키지
│   │   ├── handler/   # HTTP 핸들러
│   │   ├── service/   # 비즈니스 로직
│   │   └── k8s/       # K8s 클라이언트
│   └── go.mod
├── frontend/          # Vue 앱
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/    # Pinia stores
│   │   └── api/       # API 클라이언트
│   └── package.json
├── k8s/               # K8s 매니페스트
└── Dockerfile
```

## 기능
- [ ] 클러스터 개요 (노드 상태, 리소스 사용량)
- [ ] 네임스페이스별 파드/서비스 목록
- [ ] 서비스 상태 모니터링 (Harbor, Grafana, etc.)
- [ ] 로그 뷰어 (Loki 연동)
- [ ] 배포 관리

## 명령어
```bash
# Backend
cd backend && go run cmd/main.go

# Frontend
cd frontend && npm run dev

# Build
docker build -t homelab-dashboard .
```

## API 설계
- `GET /api/v1/cluster/nodes` - 노드 목록
- `GET /api/v1/cluster/namespaces` - 네임스페이스 목록
- `GET /api/v1/namespaces/:ns/pods` - 파드 목록
- `GET /api/v1/namespaces/:ns/services` - 서비스 목록
- `GET /api/v1/services/status` - 외부 서비스 상태

## 규칙
- 에러는 항상 JSON 형식으로 반환
- K8s API 호출은 service 레이어에서만
- 프론트엔드 상태관리는 Pinia 사용
- 컴포넌트는 Composition API + `<script setup>` 사용
