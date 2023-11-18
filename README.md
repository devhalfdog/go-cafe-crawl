# Cafe Crawler

## Project
```
/cmd : 애플리케이션 폴더
    /crawl : 과거 글을 수집하는 애플리케이션
    /watcher : 최신 글을 수집하는 애플리케이션
/config
/crawl
```

## Commit Convention

```md
feat: 새로운 기능 추가
fix: 버그 수정
docs: 문서 수정
remove: 파일 삭제할 경우
chore: 빌드 관련 파일 수정, 패키지 추가 및 수정
```

<br />

## Package Policy

```cpp
Should Only Use packages that can working both Windows and Linux.
if not, you should solve problems with team member.
```

<br />

## Git Branch Policy

```cpp
Branch Name Must be Splitted by "_" and in lowercase.
Runtime environment must be x86_64. (No ARM, 32bit)

main : Release(Deploy/Production)
develop : Debug(development)
git_testing : Can test freely git.
```

<br />

## Pull Request Policy

```cpp
Must be use template
-> path: (root directory)/.gitlab/merge_request_templates/default.md
```
