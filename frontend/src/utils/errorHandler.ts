export function handleError(error: any, context: string) {
  const message = error?.message || error?.toString() || 'Unknown error'
  console.error(`[${context}]`, error)

  // Show toast notification
  window.dispatchEvent(new CustomEvent('show-toast', {
    detail: {
      message: `${context}: ${message}`,
      type: 'error'
    }
  }))
}
