package goscraper

import (
	"strings"
)

type ContentType string

const (
	ContentTypeEcommerce   ContentType = "ecommerce"
	ContentTypeNews        ContentType = "news"
	ContentTypeBlog        ContentType = "blog"
	ContentTypeSocialMedia ContentType = "social_media"
	ContentTypeVideo       ContentType = "video"
	ContentTypeJob         ContentType = "job"
	ContentTypeRealEstate ContentType = "real_estate"
	ContentTypeRecipe      ContentType = "recipe"
	ContentTypeEvent       ContentType = "event"
	ContentTypeGeneral     ContentType = "general"
)

type ContentDetector struct {
	patterns map[ContentType][]string
	domains  map[ContentType][]string
}

func NewContentDetector() *ContentDetector {
	return &ContentDetector{
		patterns: map[ContentType][]string{
			ContentTypeEcommerce: {
				"price", "cart", "buy", "shop", "product", "store", "checkout",
				"fiyat", "sepet", "satın", "ürün", "mağaza", "alışveriş",
				"add to cart", "add-to-cart", "product-price", "buy-now",
			},
			ContentTypeNews: {
				"news", "article", "headline", "breaking", "story", "reporter",
				"haber", "makale", "başlık", "gazete", "muhabir",
				"published", "author", "byline", "news-article",
			},
			ContentTypeBlog: {
				"blog", "post", "comment", "author", "category", "tag",
				"yazı", "yorum", "kategori", "etiket",
				"blog-post", "article-content", "post-meta",
			},
			ContentTypeSocialMedia: {
				"follow", "like", "share", "tweet", "post", "profile",
				"takip", "beğen", "paylaş", "gönderi", "profil",
				"social-share", "follow-button", "user-profile",
			},
			ContentTypeVideo: {
				"video", "play", "watch", "duration", "views", "subscribe",
				"izle", "oynat", "süre", "görüntüleme", "abone",
				"video-player", "play-button", "video-duration",
			},
			ContentTypeJob: {
				"job", "career", "apply", "salary", "position", "hiring",
				"iş", "kariyer", "başvur", "maaş", "pozisyon", "işe alım",
				"job-listing", "apply-now", "job-description",
			},
			ContentTypeRealEstate: {
				"property", "rent", "sale", "bedroom", "bathroom", "sqft",
				"emlak", "kiralık", "satılık", "oda", "banyo", "metrekare",
				"property-details", "real-estate", "for-sale",
			},
			ContentTypeRecipe: {
				"recipe", "ingredient", "cooking", "preparation", "serves",
				"tarif", "malzeme", "pişirme", "hazırlık", "kişilik",
				"recipe-ingredients", "cooking-time", "prep-time",
			},
			ContentTypeEvent: {
				"event", "ticket", "date", "venue", "register", "attend",
				"etkinlik", "bilet", "tarih", "mekan", "kayıt", "katıl",
				"event-details", "buy-tickets", "event-date",
			},
		},
		domains: map[ContentType][]string{
			ContentTypeEcommerce: {
				"amazon", "ebay", "shopify", "trendyol", "hepsiburada", "n11",
				"gittigidiyor", "ciceksepeti", "morhipo", "koton", "lcwaikiki",
				"defacto", "boyner", "teknosa", "vatan", "mediamarkt",
			},
			ContentTypeNews: {
				"cnn", "bbc", "reuters", "hurriyet", "milliyet", "sabah",
				"sozcu", "cumhuriyet", "haberturk", "ntv", "cnnturk",
				"aa.com.tr", "dha.com.tr", "anha.com.tr",
			},
			ContentTypeSocialMedia: {
				"facebook", "twitter", "instagram", "linkedin", "youtube",
				"tiktok", "snapchat", "pinterest", "reddit", "discord",
			},
			ContentTypeVideo: {
				"youtube", "vimeo", "dailymotion", "twitch", "netflix",
				"bluTV", "exxen", "gain", "puhutv",
			},
			ContentTypeJob: {
				"linkedin", "indeed", "glassdoor", "kariyer.net", "yenibiris",
				"secretcv", "monster", "ziprecruiter", "jobsdb",
			},
		},
	}
}

func (cd *ContentDetector) DetectContentType(url, html string) ContentType {
	domain := extractDomainFromURL(url)
	
	for contentType, domains := range cd.domains {
		for _, d := range domains {
			if strings.Contains(strings.ToLower(domain), d) {
				return contentType
			}
		}
	}
	
	htmlLower := strings.ToLower(html)
	scores := make(map[ContentType]int)
	
	for contentType, patterns := range cd.patterns {
		for _, pattern := range patterns {
			count := strings.Count(htmlLower, pattern)
			scores[contentType] += count
		}
	}
	
	maxScore := 0
	detectedType := ContentTypeGeneral
	
	for contentType, score := range scores {
		if score > maxScore {
			maxScore = score
			detectedType = contentType
		}
	}
	
	if maxScore < 3 {
		return ContentTypeGeneral
	}
	
	return detectedType
}

func extractDomainFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		domain := parts[2]
		if strings.HasPrefix(domain, "www.") {
			domain = domain[4:]
		}
		return domain
	}
	return url
}

type SmartExtractor struct {
	detector *ContentDetector
}

func NewSmartExtractor() *SmartExtractor {
	return &SmartExtractor{
		detector: NewContentDetector(),
	}
}

func (se *SmartExtractor) ExtractSmart(resp *Response) *SmartData {
	contentType := se.detector.DetectContentType(resp.URL, resp.Body)
	parser := NewParser(resp.Document)
	
	baseData := &SmartData{
		URL:         resp.URL,
		ContentType: contentType,
		Title:       parser.ExtractTitle(),
		Description: getMetaDescription(parser),
		Images:      parser.ExtractImages(),
		Links:       parser.ExtractLinks(),
		MetaTags:    parser.ExtractMetaTags(),
	}
	
	switch contentType {
	case ContentTypeEcommerce:
		baseData.Products = se.extractProducts(parser, resp.URL)
	case ContentTypeNews:
		baseData.Article = se.extractArticle(parser)
	case ContentTypeBlog:
		baseData.BlogPost = se.extractBlogPost(parser)
	case ContentTypeJob:
		baseData.JobListing = se.extractJobListing(parser)
	case ContentTypeRealEstate:
		baseData.Property = se.extractProperty(parser)
	case ContentTypeRecipe:
		baseData.Recipe = se.extractRecipe(parser)
	case ContentTypeEvent:
		baseData.Event = se.extractEvent(parser)
	case ContentTypeVideo:
		baseData.Video = se.extractVideo(parser)
	}
	
	return baseData
}

type SmartData struct {
	URL         string      `json:"url"`
	ContentType ContentType `json:"content_type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Images      []Image     `json:"images"`
	Links       []Link      `json:"links"`
	MetaTags    map[string]string `json:"meta_tags"`
	
	Products    []SmartProduct    `json:"products,omitempty"`
	Article     *Article          `json:"article,omitempty"`
	BlogPost    *BlogPost         `json:"blog_post,omitempty"`
	JobListing  *JobListing       `json:"job_listing,omitempty"`
	Property    *Property         `json:"property,omitempty"`
	Recipe      *Recipe           `json:"recipe,omitempty"`
	Event       *Event            `json:"event,omitempty"`
	Video       *Video            `json:"video,omitempty"`
}

type SmartProduct struct {
	Name        string   `json:"name"`
	Price       string   `json:"price"`
	OriginalPrice string `json:"original_price,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	Brand       string   `json:"brand,omitempty"`
	Rating      string   `json:"rating,omitempty"`
	Reviews     string   `json:"reviews,omitempty"`
	ImageURL    string   `json:"image_url,omitempty"`
	URL         string   `json:"url,omitempty"`
	InStock     bool     `json:"in_stock"`
	Features    []string `json:"features,omitempty"`
}

type Article struct {
	Headline    string    `json:"headline"`
	Subheadline string    `json:"subheadline,omitempty"`
	Author      string    `json:"author,omitempty"`
	PublishDate string    `json:"publish_date,omitempty"`
	Content     string    `json:"content"`
	Category    string    `json:"category,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}

type BlogPost struct {
	Title       string   `json:"title"`
	Author      string   `json:"author,omitempty"`
	PublishDate string   `json:"publish_date,omitempty"`
	Content     string   `json:"content"`
	Categories  []string `json:"categories,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Comments    int      `json:"comments,omitempty"`
}

type JobListing struct {
	Title       string   `json:"title"`
	Company     string   `json:"company,omitempty"`
	Location    string   `json:"location,omitempty"`
	Salary      string   `json:"salary,omitempty"`
	JobType     string   `json:"job_type,omitempty"`
	Experience  string   `json:"experience,omitempty"`
	Description string   `json:"description"`
	Requirements []string `json:"requirements,omitempty"`
	Benefits    []string `json:"benefits,omitempty"`
	PostDate    string   `json:"post_date,omitempty"`
}

type Property struct {
	Title       string   `json:"title"`
	Price       string   `json:"price"`
	Location    string   `json:"location,omitempty"`
	PropertyType string  `json:"property_type,omitempty"`
	Bedrooms    string   `json:"bedrooms,omitempty"`
	Bathrooms   string   `json:"bathrooms,omitempty"`
	Area        string   `json:"area,omitempty"`
	Features    []string `json:"features,omitempty"`
	Description string   `json:"description"`
}

type Recipe struct {
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	PrepTime     string   `json:"prep_time,omitempty"`
	CookTime     string   `json:"cook_time,omitempty"`
	TotalTime    string   `json:"total_time,omitempty"`
	Servings     string   `json:"servings,omitempty"`
	Ingredients  []string `json:"ingredients,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
	Nutrition    string   `json:"nutrition,omitempty"`
}

type Event struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Date        string   `json:"date,omitempty"`
	Time        string   `json:"time,omitempty"`
	Venue       string   `json:"venue,omitempty"`
	Location    string   `json:"location,omitempty"`
	Price       string   `json:"price,omitempty"`
	Organizer   string   `json:"organizer,omitempty"`
	Category    string   `json:"category,omitempty"`
}

type Video struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Duration    string `json:"duration,omitempty"`
	Views       string `json:"views,omitempty"`
	Author      string `json:"author,omitempty"`
	PublishDate string `json:"publish_date,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	URL         string `json:"url,omitempty"`
}