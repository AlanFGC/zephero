import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useAuth } from "@/contexts/auth-context";
import { useDebouncedCallback } from "use-debounce";

const SignUpForm = () => {
  const authClient = useAuth();
  

  const checkUsernameDebounced = useDebouncedCallback(
    async (value) => {
      const response = await authClient.isUsernameAvailable({
      username: value})
      if (!response?.data?.available) {
        throw new Error("Username is already taken")
      } else if (response?.data?.available) {
        console.log("Username is available")
      }

      if (response?.error) {
        console.error("Error checking username availability:", response);
      }
    },
    1000,
  )

  return (
    <form>
      <div className="flex flex-col gap-6">
        <div className="grid gap-2">
          <Label htmlFor="username">Username</Label>
          <Input
            id="username"
            type="text"
            placeholder="Enter username"
            required
            onChange={(e) => checkUsernameDebounced(e.target.value)}
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="email">Email (optional)</Label>
          <Input
            id="email"
            type="email"
            placeholder="m@example.com"
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="password">Password</Label>
          <Input id="password" type="password" required />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="confirmPassword">Confirm Password</Label>
          <Input id="confirmPassword" type="password" required />
        </div>
        <Button type="submit" className="w-full">
          Sign Up
        </Button>
      </div>
    </form>
  )
}

export {
  SignUpForm,
}