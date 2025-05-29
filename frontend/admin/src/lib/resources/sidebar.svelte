<script lang="ts">
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { House, LayoutGrid, User, Box, Tags, ChartSpline, Settings} from "lucide-svelte";

    const MenuItems = [
        { name: "Home", url: "/", icon: House },
        { name: "Categories", url: "/categories", icon: LayoutGrid },
        { name: "Brands", url: "/brands", icon: Tags },
        { name: "Products", url: "/products", icon: Box },
        { name: "Customers", url: "/customers", icon: User },
        { name: "Controls", url: "/controls", icon: Settings },
        { name: "Dashboard", url: "/dashboard", icon: ChartSpline },
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
        return matchedItem ? matchedItem.name : "Home"; // Default to Home if no match
    });

    function navigateTo(name: string, url: string) {
        goto(url);
    }
</script>

<style>
    .sidebar {
        width: 60px;
        height: 100vh;
        background-color: #F8F9FA;
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 16px 0;
        position: fixed;
        left: 0;
        top: 0;
        z-index: 100;
        transition: width 0.3s;
    }

    .menu-items {
        display: flex;
        flex-direction: column;
        width: 100%;
        gap: 8px;
    }

    .menu-item {
        width: 100%;
        display: flex;
        align-items: center;
        padding: 10px 0;
        color: #6C757D;
        position: relative;
        cursor: pointer;
        transition: all 0.2s;
        text-decoration: none;
        display: flex;
    }

    .menu-item:hover {
        color: black;
    }

    .menu-item.active {
        color: black;
    }

    .menu-item.active::before {
        content: "";
        position: absolute;
        left: 0;
        top: 0;
        width: 3px;
        height: 100%;
        background-color: gray;
    }

    .icon-container {
        width: 60px;
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .tooltip {
        position: absolute;
        left: 60px;
        background-color: black;
        color: white;
        padding: 6px 12px;
        border-radius: 4px;
        white-space: nowrap;
        opacity: 0;
        visibility: hidden;
        transition: opacity 0.2s, visibility 0.2s;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
        z-index: 101;
    }

    .menu-item:hover .tooltip {
        opacity: 1;
        visibility: visible;
    }
</style>

<div class="sidebar">    
    <div class="menu-items">
        {#each MenuItems as item}
        <a
            href={item.url}
            class="menu-item"
            class:active={isActive(item.url)}
            onclick = {() => navigateTo(item.name, item.url)}
        >
            <div class="icon-container">
                <item.icon size={20} />
            </div>
            <div class="tooltip">{item.name}</div>
        </a> 
        {/each}
    </div>
</div>