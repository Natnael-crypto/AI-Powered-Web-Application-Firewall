import {useToast} from '../hooks/useToast'

function ToastExample() {
  const {addToast} = useToast()

  const showToasts = () => {
    addToast('This is an info message')
    addToast('Success!', 'success')
    addToast('Warning!', 'warning')
    addToast('Error occurred', 'error')
    addToast('Another message')
    addToast('One more toast')
  }

  return (
    <div className="p-8">
      <button
        onClick={showToasts}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        Show Multiple Toasts
      </button>
    </div>
  )
}

export default ToastExample
