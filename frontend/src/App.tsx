import { Calculator } from './calculator/Calculator'
import './App.css'

export default function App() {
  return (
    <main className="app">
      <h1>Calculator</h1>
      <p className="app__lede">Arithmetic runs on the Go API, not in the browser.</p>
      <Calculator />
    </main>
  )
}
