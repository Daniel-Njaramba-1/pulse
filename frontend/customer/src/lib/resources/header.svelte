<script lang="ts">
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { House, ShoppingCart, Heart, User, Search } from 'lucide-svelte';

    const MenuItems = [
        { name: "Home", url: "/", icon: House }, 
        { name: "Cart", url: "/cart", icon: ShoppingCart },
        { name: "Wishlist", url: "/wishlist", icon: Heart },
    ];

    const currentPathname = $derived(page.url.pathname);
    
    const isActive = $derived((url: string) => {
        if (url === "/" && currentPathname === "/") {
            return true;
        }
        return url !== "/" && currentPathname.startsWith(url);
    });

    const activeItem = $derived(() => {
        const matchedItem = MenuItems.find(item => isActive(item.url));
        return matchedItem ? matchedItem.name : "Home";
    });

    function navigateTo(name: string, url: string) {
        goto(url);
    }
</script>


<!-- ========================= HEADER COMPONENT ========================= -->
<header class="main-header">
    
    <!-- Brand Logo Section -->
    <div class="logo-container">
        <a href="/landing" class="brand-title">Pulse</a>
    </div>

    <!-- Navigation Menu Section -->
    <nav class="navigation-wrapper">
        <ul class="nav-menu">
            {#each MenuItems as item (item.name)}
                <li class="nav-item" class:active={isActive(item.url)}>
                    <a 
                        href={item.url} 
                        class="nav-link"
                        onclick={() => navigateTo(item.name, item.url)}
                    >
                        <div class="icon-wrapper">
                            <item.icon size={22} strokeWidth={1.8} />
                        </div>
                        
                        <span class="nav-label">
                            {item.name}
                        </span>
                    </a>
                </li>
            {/each}
        </ul>
    </nav>

</header>


<style>
    /* ========================= MAIN HEADER STYLES ========================= */
    .main-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        
        padding: 20px 32px;
        margin: 0 24px 32px 24px;
        
        background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
        border: 1px solid rgba(226, 232, 240, 0.6);
        border-radius: 16px;
        box-shadow: 
            0 1px 3px rgba(0, 0, 0, 0.05),
            0 4px 12px rgba(0, 0, 0, 0.02);
            
        backdrop-filter: blur(8px);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .main-header:hover {
        box-shadow: 
            0 2px 8px rgba(0, 0, 0, 0.08),
            0 8px 24px rgba(0, 0, 0, 0.04);
    }


    /* ========================= LOGO SECTION STYLES ========================= */
    .logo-container {
        display: flex;
        align-items: center;
        padding: 8px 0;
    }

    .brand-title {
        margin: 0;
        font-size: 32px;
        font-weight: 800;
        
        background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        
        letter-spacing: -0.02em;
        text-transform: uppercase;
        
        transition: all 0.3s ease;
        cursor: pointer;
    }

    .brand-title:hover {
        transform: scale(1.02);
        background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
    }


    /* ========================= NAVIGATION STYLES ========================= */
    .navigation-wrapper {
        display: flex;
        align-items: center;
    }

    .nav-menu {
        display: flex;
        list-style: none;
        margin: 0;
        padding: 0;
        
        gap: 8px;
        align-items: center;
    }

    .nav-item {
        position: relative;
        border-radius: 12px;
        overflow: hidden;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .nav-item:hover {
        transform: translateY(-2px);
    }


    /* ========================= NAVIGATION LINKS ========================= */
    .nav-link {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        
        padding: 16px 20px;
        min-width: 80px;
        
        text-decoration: none;
        color: #64748b;
        font-size: 13px;
        font-weight: 500;
        letter-spacing: 0.02em;
        
        background: transparent;
        border-radius: 12px;
        
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        overflow: hidden;
    }

    .nav-link::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        
        background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
        opacity: 0;
        
        transition: opacity 0.3s ease;
        z-index: -1;
    }

    .nav-link:hover {
        color: #1e293b;
        transform: scale(1.05);
    }

    .nav-link:hover::before {
        opacity: 1;
    }


    /* ========================= ICON WRAPPER ========================= */
    .icon-wrapper {
        display: flex;
        align-items: center;
        justify-content: center;
        
        margin-bottom: 6px;
        padding: 2px;
        
        border-radius: 8px;
        transition: all 0.3s ease;
    }

    .nav-link:hover .icon-wrapper {
        transform: scale(1.1);
    }


    /* ========================= NAVIGATION LABELS ========================= */
    .nav-label {
        font-weight: 600;
        text-transform: capitalize;
        line-height: 1.2;
        
        transition: all 0.3s ease;
    }


    /* ========================= ACTIVE STATE STYLES ========================= */
    .nav-item.active .nav-link {
        color: #0f172a;
        font-weight: 700;
        background: linear-gradient(135deg, #e2e8f0 0%, #cbd5e1 100%);
        
        box-shadow: 
            inset 0 1px 2px rgba(0, 0, 0, 0.05),
            0 2px 8px rgba(0, 0, 0, 0.1);
    }

    .nav-item.active .nav-link::before {
        opacity: 0;
    }

    .nav-item.active .icon-wrapper {
        transform: scale(1.1);
        color: #0f172a;
    }

    .nav-item.active .nav-label {
        font-weight: 800;
        text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    }


    /* ========================= RESPONSIVE DESIGN ========================= */
    @media (max-width: 768px) {
        .main-header {
            padding: 16px 20px;
            margin: 0 16px 24px 16px;
            border-radius: 12px;
        }
        
        .brand-title {
            font-size: 28px;
        }
        
        .nav-menu {
            gap: 4px;
        }
        
        .nav-link {
            padding: 12px 16px;
            min-width: 70px;
            font-size: 12px;
        }
        
        .icon-wrapper {
            margin-bottom: 4px;
        }
    }

    @media (max-width: 480px) {
        .main-header {
            padding: 12px 16px;
            margin: 0 12px 20px 12px;
            border-radius: 10px;
        }
        
        .brand-title {
            font-size: 24px;
        }
        
        .nav-link {
            padding: 10px 12px;
            min-width: 60px;
            font-size: 11px;
        }
    }
</style>