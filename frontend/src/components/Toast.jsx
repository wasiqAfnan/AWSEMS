export default function Toast({ type, message }) {
  const colors =
    type === 'success'
      ? 'bg-green-600 text-white'
      : 'bg-red-600 text-white'

  return (
    <div
      className={`fixed bottom-6 right-6 z-50 px-5 py-3 rounded-xl shadow-lg text-sm font-medium ${colors} transition-all`}
    >
      {message}
    </div>
  )
}
