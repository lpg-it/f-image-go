package fimage

import "time"

// UploadResponse represents the response from an upload operation.
type UploadResponse struct {
	// Success indicates if the upload was successful.
	Success bool `json:"success"`

	// Status is the HTTP status code.
	Status int `json:"status"`

	// Data contains the uploaded file information.
	Data *UploadData `json:"data"`
}

// UploadData contains the details of an uploaded file.
type UploadData struct {
	// ID is the unique identifier of the file.
	ID int64 `json:"id"`

	// URL is the direct URL to the original image.
	URL string `json:"url"`

	// MediumURL is the URL to the medium-sized variant (if available).
	MediumURL *string `json:"medium_url,omitempty"`

	// ThumbnailURL is the URL to the thumbnail variant (if available).
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`

	// OriginalName is the original filename.
	OriginalName string `json:"original_name"`

	// Description is the file description.
	Description string `json:"description"`

	// Size is the file size in bytes.
	Size int64 `json:"size"`

	// Width is the image width in pixels.
	Width int `json:"width"`

	// Height is the image height in pixels.
	Height int `json:"height"`

	// MimeType is the MIME type of the file.
	MimeType string `json:"mime_type"`

	// IsFlash indicates if this was a flash upload (deduplicated).
	IsFlash bool `json:"is_flash"`
}

// File represents a file in the user's library.
type File struct {
	// ID is the unique identifier of the file.
	ID int64 `json:"id"`

	// AlbumID is the ID of the album this file belongs to (if any).
	AlbumID *int64 `json:"album_id,omitempty"`

	// AlbumName is the name of the album (if any).
	AlbumName *string `json:"album_name,omitempty"`

	// OriginalName is the original filename.
	OriginalName string `json:"original_name"`

	// Description is the file description.
	Description string `json:"description"`

	// URL is the direct URL to the original image.
	URL string `json:"url"`

	// MediumURL is the URL to the medium-sized variant (if available).
	MediumURL *string `json:"medium_url,omitempty"`

	// ThumbnailURL is the URL to the thumbnail variant (if available).
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`

	// Size is the file size in bytes.
	Size int64 `json:"size"`

	// Width is the image width in pixels.
	Width int `json:"width"`

	// Height is the image height in pixels.
	Height int `json:"height"`

	// MimeType is the MIME type of the file.
	MimeType string `json:"mime_type"`

	// CreatedAt is the file creation timestamp.
	CreatedAt string `json:"created_at"`

	// DeletedAt is the soft deletion timestamp (for trash items).
	DeletedAt *string `json:"deleted_at,omitempty"`
}

// FilesListResponse represents the response from listing files.
type FilesListResponse struct {
	// Files is the list of files.
	Files []File `json:"files"`

	// Total is the total number of files.
	Total int64 `json:"total"`

	// Page is the current page number.
	Page int `json:"page"`

	// Limit is the number of items per page.
	Limit int `json:"limit"`

	// AlbumID is the album filter (if applied).
	AlbumID *int64 `json:"album_id,omitempty"`

	// Query is the search query (for search results).
	Query string `json:"query,omitempty"`
}

// Album represents an album.
type Album struct {
	// ID is the unique identifier of the album.
	ID int64 `json:"id"`

	// Name is the album name.
	Name string `json:"name"`

	// Description is the album description.
	Description string `json:"description"`

	// FileCount is the number of files in the album.
	FileCount int64 `json:"file_count"`

	// CreatedAt is the album creation timestamp.
	CreatedAt string `json:"created_at"`
}

// AlbumsListResponse represents the response from listing albums.
type AlbumsListResponse struct {
	// Albums is the list of albums.
	Albums []Album `json:"albums"`
}

// ShareLink represents a share link.
type ShareLink struct {
	// ID is the unique identifier of the share link.
	ID int64 `json:"id"`

	// Token is the share token.
	Token string `json:"token"`

	// ShareURL is the full share URL.
	ShareURL string `json:"share_url"`

	// FileID is the ID of the shared file (for file shares).
	FileID *int64 `json:"file_id,omitempty"`

	// AlbumID is the ID of the shared album (for album shares).
	AlbumID *int64 `json:"album_id,omitempty"`

	// FileName is the name of the shared file (if applicable).
	FileName *string `json:"file_name,omitempty"`

	// AlbumName is the name of the shared album (if applicable).
	AlbumName *string `json:"album_name,omitempty"`

	// HasPassword indicates if the share is password-protected.
	HasPassword bool `json:"has_password"`

	// ExpiresAt is the expiration timestamp (if set).
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// MaxViews is the maximum number of views allowed (if set).
	MaxViews *int64 `json:"max_views,omitempty"`

	// ViewCount is the current view count.
	ViewCount int64 `json:"view_count"`

	// IsActive indicates if the share link is active.
	IsActive bool `json:"is_active"`

	// CreatedAt is the share link creation timestamp.
	CreatedAt time.Time `json:"created_at"`
}

// SharesListResponse represents the response from listing share links.
type SharesListResponse struct {
	// Shares is the list of share links.
	Shares []ShareLink `json:"shares"`

	// Total is the total number of share links.
	Total int64 `json:"total"`

	// Page is the current page number.
	Page int `json:"page"`

	// Limit is the number of items per page.
	Limit int `json:"limit"`
}

// SharedContent represents the content accessed via a share link.
type SharedContent struct {
	// Type is either "file" or "album".
	Type string `json:"type"`

	// File is the shared file (for file shares).
	File *File `json:"file,omitempty"`

	// Album is the shared album (for album shares).
	Album *Album `json:"album,omitempty"`

	// Files is the list of files in the shared album.
	Files []File `json:"files,omitempty"`

	// RequiresPassword indicates if a password is required.
	RequiresPassword bool `json:"requires_password"`
}

// Tag represents a tag.
type Tag struct {
	// ID is the unique identifier of the tag.
	ID int64 `json:"id"`

	// Name is the tag name.
	Name string `json:"name"`

	// Color is the tag color (hex format).
	Color string `json:"color"`

	// FileCount is the number of files with this tag.
	FileCount int64 `json:"file_count"`
}

// TagsListResponse represents the response from listing tags.
type TagsListResponse []Tag

// TrashListResponse represents the response from listing trash items.
type TrashListResponse struct {
	// Files is the list of trashed files.
	Files []File `json:"files"`

	// Total is the total number of trashed files.
	Total int64 `json:"total"`

	// Page is the current page number.
	Page int `json:"page"`

	// Limit is the number of items per page.
	Limit int `json:"limit"`
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	// Success indicates if the operation was successful.
	Success bool `json:"success"`

	// Message is a human-readable message.
	Message string `json:"message"`

	// DeletedCount is the number of successfully deleted items.
	DeletedCount int `json:"deleted_count"`

	// FailedCount is the number of items that failed to delete.
	FailedCount int `json:"failed_count"`

	// FailedDeletions contains details about failed deletions.
	FailedDeletions []FailedDeletion `json:"failed_deletions,omitempty"`
}

// FailedDeletion represents a failed deletion with reason.
type FailedDeletion struct {
	// FileID is the ID of the file that failed to delete.
	FileID int64 `json:"file_id"`

	// FileName is the name of the file.
	FileName string `json:"file_name"`

	// Reason is why the deletion failed.
	Reason string `json:"reason"`

	// ShareLinks are the share links that blocked deletion.
	ShareLinks []ShareLink `json:"share_links,omitempty"`
}

// BatchDeleteResponse represents the response from a batch delete operation.
type BatchDeleteResponse struct {
	// Deleted is the number of successfully deleted items.
	Deleted int `json:"deleted"`

	// Failed is the number of items that failed to delete.
	Failed int `json:"failed"`

	// Message is a human-readable message.
	Message string `json:"message"`
}

// RestoreResponse represents the response from a restore operation.
type RestoreResponse struct {
	// Message is a human-readable message.
	Message string `json:"message"`

	// Restored is the number of restored files (for batch restore).
	Restored int `json:"restored,omitempty"`

	// Failed is the number of files that failed to restore.
	Failed int `json:"failed,omitempty"`
}

// MessageResponse represents a simple message response.
type MessageResponse struct {
	// Message is the response message.
	Message string `json:"message"`

	// Info provides additional information.
	Info string `json:"info,omitempty"`
}
