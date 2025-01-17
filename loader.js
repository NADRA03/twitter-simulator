export async function loadSidebar() {
    const sidebar = document.getElementById('sidebar');
    const sidebar2 = document.getElementById('sidebar2');

    try {
        const sidebarModule = await import('./sidebar.js');
        sidebar.innerHTML = sidebarModule.render();


        const sidebar2Module = await import('./sidebar2.js');
        sidebar2.innerHTML = sidebar2Module.render();
        // sidebar2Module.initialize();
    } catch (error) {
        console.error('Error loading sidebar modules:', error);
    }
}

export async function loadAChatSides() {
    const sidebar2 = document.getElementById('sidebar2');

    try {
        const chatModule = await import('./chat.js');
        sidebar2.innerHTML = chatModule.render();
    } catch (error) {
        console.error('Error loading chat sidebar module:', error);
    }
}

export async function loadPage(page) {
    const app = document.getElementById('app');

    if (page === 'log-in') {
        document.body.classList.add('login-page');
    } else {
        document.body.classList.remove('login-page');
    }

    try {
        let module;
        const parts = page.split('/');
        const basePage = parts[0];

        switch (basePage) {
            case 'home':
                module = await import('./home.js');
                module.initialize();
                break;
            case 'search':
                    module = await import('./search.js');
                    app.innerHTML = module.render();
                    module.initialize();
                    break;
            case '401':
                module = await import('./401.js');
                app.innerHTML = module.render();
                break;
            case 'sign-up':
                module = await import('./sign-up.js');
                app.innerHTML = module.render();
                break;
            case 'log-in':
                module = await import('./log-in.js');
                app.innerHTML = module.render();
                module.initialize();
                break;
            case 'post':
                const postId = parts[1];
                module = await import('./post.js');
                app.innerHTML = module.render(postId);
                module.initialize(postId);
                if (window.location.pathname !== `/post/${postId}`) {
                    history.pushState({}, '', `/post/${postId}`);
                }
                break;
            case 'chats':
                module = await import('./chats.js');
                app.innerHTML = module.render();
                module.initialize();
                break;
            case 'profile':
                    module = await import('./Profile.js');
                    module.initialize();
                    app.innerHTML = module.render();
                    break;
            case 'createPost' :
                        module = await import('./createPost.js'); 
                        app.innerHTML = module.render(); 
                        module.initializeCreatePost(); 
                        break;
            case 'comment':
                        module = await import('./comment.js'); 
                        app.innerHTML = module.render(); 
                        module.initialize(); 
                        break;
            case 'chat':
                
                module = await import('./chats.js');
                app.innerHTML = module.render();
                module.initialize();
                const chatId = parts[1];
                if (window.location.pathname !== `/chat/${chatId}`) {
                    history.pushState({}, '', `/chat/${chatId}`);
                }
                const chatModule = await import('./chat.js');
                document.getElementById('sidebar2').innerHTML = chatModule.render(chatId);
                chatModule.initialize(chatId);
                break;
            case 'a_profile':
                    const userId = parts[1];
                    module = await import('./a_profile.js');
                    app.innerHTML = module.render(userId);
                    module.initialize(userId);
    
                    // Update URL if not already set
                    if (window.location.pathname !== `/a_profile/${userId}`) {
                        history.pushState({}, '', `/a_profile/${userId}`);
                    }
                    break;
            default:
                app.innerHTML = '<h1>404 Not Found</h1>';
        }
    } catch (error) {
        console.error('Error loading page module:', error);
        app.innerHTML = '<h1>Error loading page</h1>';
    }
}

export async function getCurrentPage() {
    const path = window.location.pathname;
    const noSidebarPages = ['/sign-up', '/log-in'];

    if (!noSidebarPages.includes(path)) {
        await loadSidebar();
    }

    if (path === '/log-in') {
        document.body.classList.add('login-page');
    } else {
        document.body.classList.remove('login-page');
    }

    if (path.startsWith('/chat/')) {
        await loadAChatSides(); 
        return `chat/${path.split('/')[2]}`;  
    } else if (path === '/sign-up') {
        return 'sign-up';
    } else if (path.startsWith('/a_profile/')) {
        return `a_profile/${path.split('/')[2]}`;
    } else if (path === '/search') {
            return 'search';    
    } else if (path === '/chats') {
        return 'chats';
    }else if (path.startsWith('/post/')) {
        return `post/${path.split('/')[2]}`;
    } else if (path === '/log-in') {
        return 'log-in';
    } else if (path === '/401') {
            return '401';
    } else if (path === '/profile') {
        return 'profile';
    }else if (path === '/createPost'){
        return 'createPost';
    } else if (path === '/addComment'){
        return 'addComment';
    } else {
        return 'home';
    }
}

async function initialize() {
    const page = await getCurrentPage();
    await loadPage(page);
}

document.addEventListener("DOMContentLoaded", initialize);

