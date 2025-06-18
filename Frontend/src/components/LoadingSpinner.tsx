export default function LoadingSpinner({fullScreen = false}: {fullScreen?: boolean}) {
  return (
    <div
      className={`flex items-center justify-center ${fullScreen ? 'h-screen' : 'h-full'}`}
      data-testid="loading-spinner"
    >
      <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
    </div>
  )
}
