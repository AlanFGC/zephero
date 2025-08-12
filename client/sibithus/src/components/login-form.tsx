import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useAuth } from "@/contexts/auth-context"

const LoginForm = () => {
  const authClient = useAuth()
  const [loginInput, setLoginInput] = useState("")
  const [password, setPassword] = useState("")

  const isEmail = (input: string) => {
    return input.includes("@")
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    const isUsingEmail = isEmail(loginInput)
    
    try {
      if (isUsingEmail) {
        const { error } = await authClient.signIn.email({
          email: loginInput,
          password: password,
        })

        if (error) {
          console.error("Login failed:", error)
        }
      } else {
        const { error } = await authClient.signIn.username({
          username: loginInput,
          password: password,
        })

        if (error) {
          console.error("Login failed:", error)
        }
      }
    } catch (err) {
      console.error("Login error:", err)
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="flex flex-col gap-6">
        <div className="grid gap-2">
          <Label htmlFor="loginInput">Email or Username</Label>
          <Input
            id="loginInput"
            type="text"
            placeholder="Enter email or username"
            required
            value={loginInput}
            onChange={(e) => setLoginInput(e.target.value)}
          />
        </div>
        <div className="grid gap-2">
          <div className="flex items-center">
            <Label htmlFor="password">Password</Label>
            <button
              type="button"
              className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
            >
              Forgot your password?
            </button>
          </div>
          <Input 
            id="password" 
            type="password" 
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <Button type="submit" className="w-full">
          Login
        </Button>
      </div>
    </form>
  )
}

export { LoginForm }
