package entities

import "time"

type Torrent struct {
	ID              int              `json:"id"`    // Utopia ID
	Title           string           `json:"title"` // Title in Utopia
	File            File             `json:"file"`
	Description     string           `json:"description"`
	FileCount       int              `json:"fileCount"` // Files count
	Category        string           `json:"category"`
	Genres          []string         `json:"genres"`
	PrimaryLanguage string           `json:"primaryLanguage"`
	Files           []File           `json:"files"`
	FilesSizes      []string         `json:"filesSizes"`
	SimilarId       string           `json:"similarId"` // ID for similar torrents in Utopia
	Similar         []SimilarTorrent `json:"similar"`   // Similar torrents info
}
type File struct {
	Name  string  `json:"name"`
	Size  string  `json:"size"`
	Media *Media  `json:"media"`
	Audio []Audio `json:"audio"`
}

type Audio struct {
	Title           string `json:"title"`
	Format          string `json:"format"`
	FormatInfo      string `json:"format_info"`
	CommercialName  string `json:"commercial_name"`
	CodecID         string `json:"codec_id"`
	BitRateMode     string `json:"bit_rate_mode"`
	BitRate         string `json:"bit_rate"`
	Channels        string `json:"channels"`
	ChannelLayout   string `json:"channel_layout"`
	SamplingRate    string `json:"sampling_rate"`
	FrameRate       string `json:"frame_rate"`
	CompressionMode string `json:"compression_mode"`
	StreamSize      string `json:"stream_size"`
	Language        string `json:"language"`
	Default         bool   `json:"default"`
	Forced          bool   `json:"forced"`
	Studio          string `json:"studio"`
}

type Media struct {
	Format             string `json:"format"`
	FormatVersion      string `json:"format_version"`
	OverallBitRateMode string `json:"overall_bit_rate_mode"`
	OverallBitRate     string `json:"overall_bit_rate"`
	FrameRate          string `json:"frame_rate"`
	Resolution         string `json:"resolution"`
}

type SimilarTorrent struct {
	ID              int       `json:"id"`              // Идентификатор :contentReference[oaicite:13]{index=13}
	Title           string    `json:"title"`           // Название :contentReference[oaicite:14]{index=14}
	URL             string    `json:"url"`             // Ссылка на страницу
	Type            string    `json:"type"`            // Тип (WEB-DL, Remux…) :contentReference[oaicite:15]{index=15}
	SizeBytes       int64     `json:"sizeBytes"`       // Размер в байтах :contentReference[oaicite:16]{index=16}
	Seeders         int       `json:"seeders"`         // Сидеры :contentReference[oaicite:17]{index=17}
	Leechers        int       `json:"leechers"`        // Личеры :contentReference[oaicite:18]{index=18}
	Completed       int       `json:"completed"`       // Завершено скачиваний :contentReference[oaicite:19]{index=19}
	UploadedAt      time.Time `json:"uploadedAt"`      // Дата добавления :contentReference[oaicite:20]{index=20}
	Thanks          int       `json:"thanks"`          // Подяк отримано :contentReference[oaicite:21]{index=21}
	Comments        int       `json:"comments"`        // Число комментариев :contentReference[oaicite:22]{index=22}
	PersonalRelease bool      `json:"personalRelease"` // Личный релиз :contentReference[oaicite:23]{index=23}
	FreeLeech       bool      `json:"freeLeech"`       // Глобальный фрилейч :contentReference[oaicite:24]{index=24}
	Attributes      []string  `json:"attributes"`      // Доп. иконки (Дубляж, Субтитры…) :contentReference[oaicite:25]{index=25}
}
