import { useState, useEffect } from 'react';
import { Download, Link, X, Save, Moon, Sun, Pause, Play, Info, XCircle, Trash2, LayoutGrid, List } from 'lucide-react';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Card, CardContent } from '@/components/ui/card';


interface Download {
  id: number;
  url: string;
  progress: number;
  status: string;
  error?: string;
  createdAt: string;
  filename: string;
  size: string;
  thumbnail: string;
  type: string;
  quality: string;
  duration: string;
  speed: string;
  timeRemaining: string;
}



const VideoDownloader = () => {
  const [url, setUrl] = useState('');
  const [downloads, setDownloads] = useState<Download[]>([]);
  const [darkMode, setDarkMode] = useState(false);
  const [selectedDownload, setSelectedDownload] = useState<Download | null>(null);
  const [viewMode, setViewMode] = useState('grid');
  const [isValidUrl, setIsValidUrl] = useState(false);

  // Mock data for initial downloads
  const mockDownloads = [
    {
      id: 1,
      url: 'https://www.youtube.com/watch?v=OlbOOzX_fDs&pp=ygUcbmF0IGdlbyB3aWxkbGlmZSBkb2N1bWVudGFyeQ%3D%3D',
      progress: 100,
      status: 'completed',
      filename: 'Nature-Documentary-HD.mp4',
      size: '1.2 GB',
      thumbnail: 'https://i.ytimg.com/vi_webp/wDhOcRa6BK0/maxresdefault.webp',
      type: 'MP4',
      quality: '1080p',
      duration: '1:23:45',
      speed: '3.2 MB/s',
      timeRemaining: '12:30'
    },
    {
      id: 2,
      url: 'https://www.youtube.com/watch?v=OlbOOzX_fDs&pp=ygUcbmF0IGdlbyB3aWxkbGlmZSBkb2N1bWVudGFyeQ%3D%3D',
      progress: 100,
      status: 'completed',
      filename: 'Cooking-Masterclass.mp4',
      size: '845 MB',
      thumbnail: 'https://i.ytimg.com/vi_webp/HUZiXcqHjw8/maxresdefault.webp',
      type: 'MP4',
      quality: '720p',
      duration: '45:12',
      speed: '3.2 MB/s',
      timeRemaining: '12:30'
    },
    {
      id: 3,
      url: 'https://www.youtube.com/watch?v=ADc8TaxAVMU&pp=ygUQbmF0IGdlbyB3aWxkbGlmZSBkb2N1bWVudGFyeQ%3D%3D',
      progress: 35,
      status: 'downloading',
      filename: 'Full-Body-Workout.mp4',
      size: '650 MB',
      thumbnail: 'https://i.ytimg.com/vi_webp/ADc8TaxAVMU/maxresdefault.webp',
      type: 'MP4',
      quality: '1080p',
      duration: '32:15',
      speed: '2.1 MB/s',
      timeRemaining: '5:45'
    }
  ];

  
  const validateUrl = (url: string) => {
    const videoUrlRegex = /^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be|vimeo\.com|tiktok\.com)\/.+$/;
    return videoUrlRegex.test(url);
  };

  
  // Handle URL input change
  const handleUrlChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newUrl = e.target.value;
    setUrl(newUrl);
    setIsValidUrl(validateUrl(newUrl));
  };

  useEffect(() => {
    setDownloads(mockDownloads);
    
    // Simulate progress for downloading video
    const interval = setInterval(() => {
      setDownloads(prevDownloads => 
        prevDownloads.map(download => {
          if (download.status === 'downloading' && download.progress < 100) {
            return {
              ...download,
              progress: Math.min(download.progress + 1, 100),
              status: download.progress + 1 === 100 ? 'completed' : 'downloading'
            };
          }
          return download;
        })
      );
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const startDownload = () => {
    if (!url) return;
    
    const newDownload = {
      id: Date.now(),
      url,
      progress: 0,
      status: 'downloading',
      filename: `video-${Date.now()}.mp4`,
      size: '750 MB',
      thumbnail: '/api/placeholder/320/180',
      type: 'MP4',
      quality: '1080p',
      duration: '25:30',
      speed: '3.2 MB/s',
      timeRemaining: '12:30'
    };
    
    setDownloads([newDownload, ...downloads]);
    setUrl('');
  };

  const togglePause = (id: number) => {
    setDownloads(downloads.map(download => {
      if (download.id === id) {
        return {
          ...download,
          status: download.status === 'paused' ? 'downloading' : 'paused'
        };
      }
      return download;
    }));
  };

  const cancelDownload = (id: number) => {
    setDownloads(downloads.map(download => {
      if (download.id === id) {
        return {
          ...download,
          status: 'cancelled',
          progress: 0
        };
      }
      return download;
    }));
  };

  
  const cancelAllDownloads = () => {
    setDownloads(downloads.map(download => ({
      ...download,
      status: download.status === 'downloading' ? 'cancelled' : download.status,
      progress: download.status === 'downloading' ? 0 : download.progress
    })));
  };

  const removeDownload = (id: number) => {
    setDownloads(downloads.filter(download => download.id !== id));
  };

  return (
    <div className={`min-h-screen ${darkMode ? 'dark bg-gray-900' : 'bg-gray-50'}`}>
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-900 dark:to-purple-900">
        <div className="container mx-auto px-4 py-12 text-center text-white">
          <h1 className="text-4xl md:text-5xl font-bold mb-4">
            Download Videos Without Watermarks
          </h1>
          <p className="text-lg md:text-xl opacity-90 max-w-2xl mx-auto">
            Download high-quality videos from YouTube, TikTok, and more without any watermarks. 
            Fast, free, and secure downloads in multiple formats.
          </p>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Header Controls */}
        <div className="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
          <div className="flex items-center gap-4">
            <button
              onClick={() => setDarkMode(!darkMode)}
              className="p-2 rounded-full hover:bg-white/10 dark:hover:bg-black/10 transition-colors"
            >
              {darkMode ? <Sun className="text-white" /> : <Moon className="text-gray-800" />}
            </button>
            <div className="flex items-center gap-2 bg-white dark:bg-gray-800 rounded-lg p-1 shadow-sm">
              <button
                onClick={() => setViewMode('grid')}
                className={`p-2 rounded-md transition-colors ${
                  viewMode === 'grid' 
                    ? 'bg-blue-500 text-white' 
                    : 'text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700'
                }`}
              >
                <LayoutGrid size={20} />
              </button>
              <button
                onClick={() => setViewMode('list')}
                className={`p-2 rounded-md transition-colors ${
                  viewMode === 'list' 
                    ? 'bg-blue-500 text-white' 
                    : 'text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700'
                }`}
              >
                <List size={20} />
              </button>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <span className="text-sm text-gray-600 dark:text-gray-300">
              {downloads.filter(d => d.status === 'completed').length} completed • 
              {downloads.filter(d => d.status === 'downloading').length} active
            </span>
            {downloads.some(d => d.status === 'downloading') && (
              <button
                onClick={cancelAllDownloads}
                className="flex items-center gap-2 px-4 py-2 text-sm text-red-500 hover:text-red-600 
                         bg-red-50 dark:bg-red-900/20 rounded-lg transition-colors"
              >
                <Trash2 size={16} />
                Cancel All Downloads
              </button>
            )}
          </div>
        </div>
        
        {/* URL Input Card */}
        <Card className="mb-8 dark:bg-gray-800/50 backdrop-blur-sm border-t border-white/20">
          <CardContent className="pt-6">
            <div className="flex flex-col md:flex-row gap-4">
              <div className="flex-1">
                <div className="relative">
                  <Link className="absolute left-3 top-3 text-gray-400" size={20} />
                  <input
                    type="text"
                    value={url}
                    onChange={handleUrlChange}
                    placeholder="Paste video URL from YouTube, TikTok, or Vimeo..."
                    className="w-full pl-10 pr-4 py-3 rounded-lg border dark:border-gray-600 
                             dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-blue-500 
                             focus:border-transparent transition-all"
                  />
                  {url && (
                    <div className={`absolute right-3 top-3 text-sm ${
                      isValidUrl ? 'text-green-500' : 'text-red-500'
                    }`}>
                      {isValidUrl ? 'Valid URL' : 'Invalid URL'}
                    </div>
                  )}
                </div>
              </div>
              <button
                onClick={startDownload}
                disabled={!isValidUrl}
                className="px-6 py-3 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg 
                         hover:from-blue-600 hover:to-blue-700 disabled:from-gray-300 disabled:to-gray-400 
                         disabled:cursor-not-allowed flex items-center gap-2 transition-all transform 
                         hover:scale-105 active:scale-95 shadow-lg hover:shadow-xl"
              >
                <Download size={20} />
                Download Video
              </button>
            </div>
          </CardContent>
        </Card>
        
        {/* Downloads Grid */}
        

        {viewMode === 'grid' ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {downloads.map(download => (
            <Card
              key={download.id}
              className="overflow-hidden dark:bg-gray-800 hover:shadow-lg transition-shadow"
            >
              <div className="relative">
                <img
                  src={download.thumbnail}
                  alt={download.filename}
                  className="w-full h-48 object-cover"
                />
                <div className="absolute top-2 right-2 flex gap-2">
                  <button
                    onClick={() => setSelectedDownload(download)}
                    className="p-1 rounded-full bg-gray-900/70 text-white hover:bg-gray-900/90"
                  >
                    <Info size={16} />
                  </button>
                  <button
                    onClick={() => removeDownload(download.id)}
                    className="p-1 rounded-full bg-gray-900/70 text-white hover:bg-gray-900/90"
                  >
                    <X size={16} />
                  </button>
                </div>
              </div>

              <CardContent className="p-4">
                <div className="mb-2">
                  <h3 className="font-medium truncate dark:text-white">
                    {download.filename}
                  </h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 truncate">
                    {download.url}
                  </p>
                </div>

                {/* Progress Section */}
                <div className="space-y-2">
                  <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                    <div
                      className={`h-full transition-all duration-300 ${
                        download.status === 'completed' ? 'bg-green-500' :
                        download.status === 'downloading' ? 'bg-blue-500' :
                        download.status === 'paused' ? 'bg-yellow-500' :
                        'bg-red-500'
                      }`}
                      style={{ width: `${download.progress}%` }}
                    />
                  </div>

                  <div className="flex justify-between items-center text-sm">
                    <span className="text-gray-600 dark:text-gray-300">
                      {download.status === 'downloading' && (
                        <>
                          {download.speed} • {download.timeRemaining} left
                        </>
                      )}
                      {download.status === 'completed' && 'Completed'}
                      {download.status === 'paused' && 'Paused'}
                      {download.status === 'cancelled' && 'Cancelled'}
                    </span>
                    <span className="text-gray-600 dark:text-gray-300">
                      {download.progress}%
                    </span>
                  </div>

                  {/* Action Buttons */}
                  <div className="flex justify-between items-center pt-2">
                    {download.status === 'completed' ? (
                      <button className="flex items-center gap-1 text-sm text-blue-500 
                                     hover:text-blue-600 dark:text-blue-400">
                        <Save size={16} />
                        Save to device
                      </button>
                    ) : (
                      <div className="flex gap-2">
                        {download.status !== 'cancelled' && (
                          <>
                            <button
                              onClick={() => togglePause(download.id)}
                              className="flex items-center gap-1 text-sm text-yellow-500 
                                       hover:text-yellow-600"
                            >
                              {download.status === 'paused' ? (
                                <Play size={16} />
                              ) : (
                                <Pause size={16} />
                              )}
                              {download.status === 'paused' ? 'Resume' : 'Pause'}
                            </button>
                            <button
                              onClick={() => cancelDownload(download.id)}
                              className="flex items-center gap-1 text-sm text-red-500 
                                       hover:text-red-600"
                            >
                              <XCircle size={16} />
                              Cancel
                            </button>
                          </>
                        )}
                      </div>
                    )}
                    <div className="text-sm text-gray-500 dark:text-gray-400">
                      {download.size} • {download.quality}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
        ) : (
          <div className="space-y-4">
            {downloads.map(download => (
              <Card
                key={download.id}
                className="overflow-hidden dark:bg-gray-800/50 backdrop-blur-sm hover:shadow-lg 
                         transition-all hover:scale-[1.01] border-t border-white/20"
              >
                <CardContent className="p-4">
                  <div className="flex gap-4">
                    <div className="w-48 h-32 flex-shrink-0">
                      <img
                        src={download.thumbnail}
                        alt={download.filename}
                        className="w-full h-full object-cover rounded-lg"
                      />
                    </div>
                    <div className="flex-1">
                      <div className="flex justify-between items-start">
                        <div>
                          <h3 className="font-medium dark:text-white truncate">
                            {download.filename}
                          </h3>
                          <p className="text-sm text-gray-500 dark:text-gray-400 truncate">
                            {download.url}
                          </p>
                        </div>
                        <div className="flex gap-2">
                          <button
                            onClick={() => setSelectedDownload(download)}
                            className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 
                                     transition-colors"
                          >
                            <Info size={16} />
                          </button>
                          <button
                            onClick={() => removeDownload(download.id)}
                            className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 
                                     transition-colors"
                          >
                            <X size={16} />
                          </button>
                        </div>
                      </div>
                      
                      {/* Progress Section */}
                      <div className="mt-4 space-y-2">
                        <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                          <div
                            className={`h-full transition-all duration-300 ${
                              download.status === 'completed' ? 'bg-gradient-to-r from-green-400 to-green-500' :
                              download.status === 'downloading' ? 'bg-gradient-to-r from-blue-400 to-blue-500' :
                              download.status === 'paused' ? 'bg-gradient-to-r from-yellow-400 to-yellow-500' :
                              'bg-gradient-to-r from-red-400 to-red-500'
                            }`}
                            style={{ width: `${download.progress}%` }}
                          />
                        </div>

                        {/* Status and Controls */}
                        <div className="flex justify-between items-center">
                          <div className="flex items-center gap-4">
                            {download.status === 'completed' ? (
                              <button className="flex items-center gap-1 text-sm text-blue-500 
                                             hover:text-blue-600 transition-colors">
                                <Save size={16} />
                                Save to device
                              </button>
                            ) : (
                              <div className="flex gap-2">
                                {download.status !== 'cancelled' && (
                                  <>
                                    <button
                                      onClick={() => togglePause(download.id)}
                                      className="flex items-center gap-1 text-sm text-yellow-500 
                                               hover:text-yellow-600 transition-colors"
                                    >
                                      {download.status === 'paused' ? (
                                        <Play size={16} />
                                      ) : (
                                        <Pause size={16} />
                                      )}
                                      {download.status === 'paused' ? 'Resume' : 'Pause'}
                                    </button>
                                    <button
                                      onClick={() => cancelDownload(download.id)}
                                      className="flex items-center gap-1 text-sm text-red-500 
                                               hover:text-red-600 transition-colors"
                                    >
                                      <XCircle size={16} />
                                      Cancel
                                    </button>
                                  </>
                                )}
                              </div>
                            )}
                          </div>
                          <div className="text-sm text-gray-500 dark:text-gray-400">
                            {download.size} • {download.quality}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        
        {/* Empty State */}
        {downloads.length === 0 && (
          <Alert className="dark:bg-gray-800/50 backdrop-blur-sm">
            <AlertDescription>
              No downloads yet. Paste a video URL above to get started.
            </AlertDescription>
          </Alert>
        )}

        {/* Download Details Modal */}
        {selectedDownload && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4">
            <Card className="w-full max-w-lg dark:bg-gray-800">
              <CardContent className="p-6">
                <div className="flex justify-between items-start mb-4">
                  <h3 className="text-xl font-bold dark:text-white">Download Details</h3>
                  <button
                    onClick={() => setSelectedDownload(null)}
                    className="text-gray-500 hover:text-gray-700 dark:text-gray-400 
                             dark:hover:text-gray-200"
                  >
                    <X size={20} />
                  </button>
                </div>
                <div className="space-y-3">
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <span className="text-gray-500 dark:text-gray-400">Filename:</span>
                    <span className="dark:text-white">{selectedDownload.filename}</span>
                    <span className="text-gray-500 dark:text-gray-400">Size:</span>
                    <span className="dark:text-white">{selectedDownload.size}</span>
                    <span className="text-gray-500 dark:text-gray-400">Quality:</span>
                    <span className="dark:text-white">{selectedDownload.quality}</span>
                    <span className="text-gray-500 dark:text-gray-400">Duration:</span>
                    <span className="dark:text-white">{selectedDownload.duration}</span>
                    <span className="text-gray-500 dark:text-gray-400">Type:</span>
                    <span className="dark:text-white">{selectedDownload.type}</span>
                    <span className="text-gray-500 dark:text-gray-400">Status:</span>
                    <span className="dark:text-white capitalize">{selectedDownload.status}</span>
                  </div>
                  <div className="pt-4">
                    <p className="text-sm text-gray-500 dark:text-gray-400 break-all">
                      Source: {selectedDownload.url}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
};

export default VideoDownloader;