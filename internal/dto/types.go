package dto

// --- JSON Structure ---------------------------

type Bookshelf struct {
	Books       []Book       `json:"books"`
	Collections []Collection `json:"collections"`
	Insights    Insights     `json:"insights"`
}

type Book struct {
	Id        string   `json:"id"`
	Isbn      string   `json:"isbn"`
	Title     string   `json:"title"`
	Subtitle  string   `json:"subtitle"`
	Authors   []string `json:"authors"`
	Year      int      `json:"year"`
	Language  string   `json:"language"`
	Pages     int      `json:"pages"`
	Genre     string   `json:"genre"`
	Tags      []string `json:"tags"`
	Cover     string   `json:"cover"`
	Link      string   `json:"link"`
	DateAdded string   `json:"date_added"`
	Status    string   `json:"status"`
	Rank      int      `json:"rank"`
	Progress  Progress `json:"progress"`
	Rating    float64  `json:"rating"`
	Review    []string `json:"review"`
	Quotes    []string `json:"quotes"`
}

type Progress struct {
	DateStarted  string `json:"date_started"`
	DateFinished string `json:"date_finished"`
	ReadingTime  string `json:"reading_time"`
	PagesRead    int    `json:"pages_read"`
}

type Collection struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Books       []string `json:"books"`
}

type Insights struct {
	FavoriteQuote Quote `json:"favoriteQuote"`
}

type Quote struct {
	Quote     string   `json:"quote"`
	Authors   []string `json:"authors"`
	BookTitle string   `json:"title"`
	Id        string   `json:"id"`
}

// ----------------------------------------------

type Stats struct {
	TotalBooks            int
	BooksFinished         int
	BooksFinishedThisYear int
	PagesRead             int
	PagesReadThisYear     int
	AverageRating         float64
	AveragePages          float64
	TopGenres             []StatCount
	BooksByStatus         []StatCount
	BooksByLanguage       []StatCount
}

type StatCount struct {
	Value string
	Count int
}

type ResolvedCollection struct {
	Name        string
	Description string
	Books       []Book
}
