package crawl

type Post struct {
	Result struct {
		ErrorCode string     `json:"errorCode,omitempty"` // 에러 발생 시 에러 코드
		Reason    string     `json:"reason,omitempty"`    // 에러 발생 시 원인
		Heads     []struct { // 말머리 리스트
			HeadID int64  `json:"headId,omitempty"` // 말머리 ID
			Head   string `json:"head,omitempty"`   // 말머리 제목
		} `json:"heads,omitempty"`
		Article struct { // 글 정보
			ID           int64    `json:"id,omitempty"`           // 글 ID
			RefArticleID int64    `json:"refArticleId,omitempty"` // 답글 ID
			Subject      string   `json:"subject,omitempty"`      // 글 제목
			Head         string   `json:"head,omitempty"`         // 말머리
			Content      string   `json:"contentHtml,omitempty"`  // 게시물 HTML
			WriteDate    int64    `json:"writeDate,omitempty"`
			Menu         struct { // 게시판 정보
				Name string `json:"name,omitempty"` // 게시판 이름
			} `json:"menu,omitempty"`
			Writer struct { // 글 작성자 정보
				ID    string `json:"id,omitempty"`          // 작성자 ID
				Nick  string `json:"nick,omitempty"`        // 작성자 닉네임
				Level int64  `json:"memberLevel,omitempty"` // 작성자 레벨
			} `json:"writer,omitempty"`
			Comments struct {
				Items []struct {
					ID      int64    `json:"id,omitempty"`      // 댓글 ID
					RefID   int64    `json:"refId,omitempty"`   // 대댓글 ID
					Content string   `json:"content,omitempty"` // 댓글 내용
					Writer  struct { // 댓글 작성자 정보
						ID   string `json:"id,omitempty"`   // 댓글 작성자 ID
						Nick string `json:"nick,omitempty"` // 댓글 작성자 닉네임
					} `json:"writer,omitempty"`
				} `json:"items,omitempty"`
			} `json:"comments,omitempty"`
		} `json:"article,omitempty"`
	} `json:"result,omitempty"`
}

/* Naver New Post JSON Struct */
type newPost struct {
	Message struct {
		Status string `json:"status"`
		Error  struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		Result struct {
			ArticleList []struct {
				ID       int64  `json:"articleId"` // 글 ID
				MenuName string `json:"menuName"`  // 메뉴 이름
				Subject  string `json:"subject"`   // 글 제목
				HeadName string `json:"headName"`  // 말머리
			} `json:"articleList"`
		} `json:"result"`
	} `json:"message"`
}

/* Naver video information json struct */
type video struct {
	Video struct {
		List []struct {
			Source string `json:"source"`
		} `json:"list"`
	} `json:"videos"`
}
