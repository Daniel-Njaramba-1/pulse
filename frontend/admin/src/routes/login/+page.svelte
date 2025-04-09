<script lang="ts">
    import { goto } from "$app/navigation";
    import { Button } from "$lib/components/ui/button/index";
    import * as Card from "$lib/components/ui/card/index";
    import * as Alert from "$lib/components/ui/alert/index";
    import { Input } from "$lib/components/ui/input/index";
    import { Label } from "$lib/components/ui/label/index";
    import { authHelpers } from "$lib/stores/auth";
    
    let username = $state<string>("");
    let password = $state<string>("");
    let loginLoading = $state<boolean>(false);
    let loginError = $state<string>("");

    async function handleLogin(event:SubmitEvent): Promise<void> {
        event.preventDefault();
        loginLoading = true;
        loginError = "";

        try {
            const result:any = await authHelpers.login(username, password);
            if (result.success) {
                goto("/");
            } else {
                loginError = result.error;
            }
        } catch (error) {
            loginError = "Error occurred while logging in";
        } finally {
            loginLoading= false;
        }
    }
</script>

<div class="flex justify-center items-center min-h-screen bg-gray-100">
    <Card.Root class="w-full max-w-md shadow-lg">
        <Card.Header class="space-y-1">
            <Card.Title class="text-2xl text-center">Login</Card.Title>
            <Card.Description class="text-center text-gray-500">
                Enter your credentials to access your account
            </Card.Description>
        </Card.Header>
        <Card.Content>
            <form onsubmit={handleLogin} class="space-y-4">
                <div class="space-y-2">
                    <Label for="username" class="text-sm font-medium">Username</Label>
                    <Input 
                        id="username"
                        type="text" 
                        bind:value={username} 
                        placeholder="Enter your username"
                        required
                        class="w-full"
                    />
                </div>
                <div class="space-y-2">
                    <Label for="password" class="text-sm font-medium">Password</Label>
                    <Input 
                        id="password"
                        type="password" 
                        bind:value={password} 
                        placeholder="Enter your password"
                        required
                        class="w-full"
                    />
                </div>
                {#if loginError}
                    <Alert.Root variant="destructive" class="my-4">
                        <Alert.Title>Authentication Error</Alert.Title> 
                        <Alert.Description>{loginError}</Alert.Description>
                    </Alert.Root>
                {/if}
                <Button 
                    type="submit" 
                    class="w-full" 
                    variant="default"
                    disabled={loginLoading}
                >
                    {loginLoading ? 'Logging in...' : 'Login'}
                </Button>
            </form>
        </Card.Content>
        <Card.Footer class="flex justify-center">
            <p class="text-sm text-gray-500">
                Don't have an account? <a href="/register" class="text-blue-600 hover:underline">Register</a>
            </p>
        </Card.Footer>
    </Card.Root>
</div>