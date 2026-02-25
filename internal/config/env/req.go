// =============================================================================
// FILE: internal/config/env/req.go
// PURPOSE: Defines all HTTP request-related configuration defaults including
//          timeouts, connection pool sizes, semaphore limits, retry parameters,
//          and worker thread counts. Ports Python of_env/values/req/req.py.
// =============================================================================

package env

import "time"

// ---------------------------------------------------------------------------
// Connection timeouts
// ---------------------------------------------------------------------------

// ConnectTimeout returns the HTTP connection timeout.
func ConnectTimeout() time.Duration {
	return GetDuration("OF_CONNECT_TIMEOUT", 35)
}

// PoolConnectTimeout returns the connection pool acquire timeout.
func PoolConnectTimeout() time.Duration {
	return GetDuration("OF_POOL_CONNECT_TIMEOUT", 60)
}

// TotalTimeout returns the total request timeout (0 = unlimited).
func TotalTimeout() time.Duration {
	return GetDuration("OF_TOTAL_TIMEOUT", 0)
}

// KeepAlive returns the HTTP keep-alive duration.
func KeepAlive() time.Duration {
	return GetDuration("OF_KEEP_ALIVE", 20)
}

// KeepAliveExpiry returns the keep-alive connection expiry.
func KeepAliveExpiry() time.Duration {
	return GetDuration("OF_KEEP_ALIVE_EXP", 10)
}

// ChunkTimeoutSec returns the timeout for individual chunk downloads.
func ChunkTimeoutSec() time.Duration {
	return GetDuration("OF_CHUNK_TIMEOUT_SEC", 120)
}

// ---------------------------------------------------------------------------
// Connection pool sizes
// ---------------------------------------------------------------------------

// MaxConnections returns the maximum total HTTP connections.
func MaxConnections() int {
	return GetInt("OF_MAX_CONNECTIONS", 200)
}

// APIMaxConnections returns the maximum connections for API requests.
func APIMaxConnections() int {
	return GetInt("OF_API_MAX_CONNECTIONS", 100)
}

// ---------------------------------------------------------------------------
// Chunk sizes (bytes)
// ---------------------------------------------------------------------------

// MaxChunkSize returns the maximum download chunk size in bytes (128 MB).
func MaxChunkSize() int64 {
	return GetInt64("OF_MAX_CHUNK_SIZE", 134217728)
}

// MinChunkSize returns the minimum download chunk size in bytes (64 KB).
func MinChunkSize() int64 {
	return GetInt64("OF_MIN_CHUNK_SIZE", 65536)
}

// ChunkUpdateCount returns the number of chunks before recalculating chunk size.
func ChunkUpdateCount() int {
	return GetInt("OF_CHUNK_UPDATE_COUNT", 12)
}

// ChunkSizeUpdateCount returns the frequency of chunk size updates.
func ChunkSizeUpdateCount() int {
	return GetInt("OF_CHUNK_SIZE_UPDATE_COUNT", 15)
}

// ChunkMemorySplit returns the memory split threshold for chunked downloads (64 MB).
func ChunkMemorySplit() int64 {
	return GetInt64("OF_CHUNK_MEMORY_SPLIT", 67108864)
}

// ChunkFileSplit returns the file split threshold for chunked downloads (64 MB).
func ChunkFileSplit() int64 {
	return GetInt64("OF_CHUNK_FILE_SPLIT", 67108864)
}

// MaxReadSize returns the maximum single read size in bytes (16 MB).
func MaxReadSize() int64 {
	return GetInt64("OF_MAX_READ_SIZE", 16777216)
}

// ---------------------------------------------------------------------------
// Semaphore limits
// ---------------------------------------------------------------------------

// ReqSemaphoreMulti returns the multiplier for request semaphores.
func ReqSemaphoreMulti() int {
	return GetInt("OF_REQ_SEMAPHORE_MULTI", 5)
}

// ScrapePaidSems returns the semaphore limit for paid content scraping.
func ScrapePaidSems() int {
	return GetInt("OF_SCRAPE_PAID_SEMS", 10)
}

// SubscriptionSems returns the semaphore limit for subscription fetching.
func SubscriptionSems() int {
	return GetInt("OF_SUBSCRIPTION_SEMS", 5)
}

// LikeMaxSems returns the maximum semaphore count for like operations.
func LikeMaxSems() int {
	return GetInt("OF_LIKE_MAX_SEMS", 12)
}

// MaxSemsBatchDownload returns the semaphore limit for batch downloads.
func MaxSemsBatchDownload() int {
	return GetInt("OF_MAX_SEMS_BATCH_DOWNLOAD", 12)
}

// MaxSemsSingleThreadDownload returns the semaphore limit for single-thread downloads.
func MaxSemsSingleThreadDownload() int {
	return GetInt("OF_MAX_SEMS_SINGLE_THREAD_DOWNLOAD", 50)
}

// SessionManagerSyncSem returns the sync session manager semaphore count.
func SessionManagerSyncSem() int {
	return GetInt("OF_SESSION_MANAGER_SYNC_SEM", 3)
}

// SessionManagerSem returns the async session manager semaphore count.
func SessionManagerSem() int {
	return GetInt("OF_SESSION_MANAGER_SEM", 10)
}

// MaxThreadWorkers returns the maximum number of thread pool workers.
func MaxThreadWorkers() int {
	return GetInt("OF_MAX_THREAD_WORKERS", 20)
}

// ---------------------------------------------------------------------------
// Session wait/retry parameters
// ---------------------------------------------------------------------------

// OFMinWaitSession returns the minimum wait between session requests.
func OFMinWaitSession() float64 {
	return GetFloat64("OF_MIN_WAIT_SESSION", 2.0)
}

// OFMaxWaitSession returns the maximum wait between session requests.
func OFMaxWaitSession() float64 {
	return GetFloat64("OF_MAX_WAIT_SESSION", 6.0)
}

// OFMinWaitExponentialSession returns the minimum exponential backoff wait.
func OFMinWaitExponentialSession() float64 {
	return GetFloat64("OF_MIN_WAIT_EXP_SESSION", 16.0)
}

// OFMaxWaitExponentialSession returns the maximum exponential backoff wait.
func OFMaxWaitExponentialSession() float64 {
	return GetFloat64("OF_MAX_WAIT_EXP_SESSION", 128.0)
}

// OFNumRetriesSession returns the maximum number of session retries.
func OFNumRetriesSession() int {
	return GetInt("OF_NUM_RETRIES_SESSION", 10)
}

// OFMinWaitAPI returns the minimum wait between API requests.
func OFMinWaitAPI() float64 {
	return GetFloat64("OF_MIN_WAIT_API", 2.0)
}

// OFMaxWaitAPI returns the maximum wait between API requests.
func OFMaxWaitAPI() float64 {
	return GetFloat64("OF_MAX_WAIT_API", 6.0)
}

// OFAuthMinWait returns the minimum wait for auth requests.
func OFAuthMinWait() float64 {
	return GetFloat64("OF_AUTH_MIN_WAIT", 3.0)
}

// OFAuthMaxWait returns the maximum wait for auth requests.
func OFAuthMaxWait() float64 {
	return GetFloat64("OF_AUTH_MAX_WAIT", 10.0)
}

// DownloadNumTriesReq returns the retry count for download requests.
func DownloadNumTriesReq() int {
	return GetInt("OF_DOWNLOAD_NUM_TRIES_REQ", 5)
}

// DownloadNumTriesCheckReq returns the retry count for download check requests.
func DownloadNumTriesCheckReq() int {
	return GetInt("OF_DOWNLOAD_NUM_TRIES_CHECK_REQ", 2)
}

// AuthNumTries returns the retry count for auth requests.
func AuthNumTries() int {
	return GetInt("OF_AUTH_NUM_TRIES", 3)
}

// MessageSleep returns the sleep duration between message fetches.
func MessageSleep() time.Duration {
	return GetDuration("OF_MESSAGE_SLEEP", 0)
}

// ---------------------------------------------------------------------------
// Proxy configuration
// ---------------------------------------------------------------------------

// Proxy returns the HTTP proxy URL (empty = no proxy).
func Proxy() string {
	return GetString("OF_PROXY", "")
}

// ProxyAuth returns the proxy authentication string (empty = no auth).
func ProxyAuth() string {
	return GetString("OF_PROXY_AUTH", "")
}
