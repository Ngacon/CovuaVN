// Mảng lưu trữ bài đăng
let posts = [];

// Hàm thêm bài đăng
function addPost() {
    const title = document.getElementById('postTitle').value.trim();
    const content = document.getElementById('postContent').value.trim();

    if (!title || !content) {
        alert("Vui lòng nhập tiêu đề và nội dung!");
        return;
    }

    // Thêm bài đăng vào mảng
    posts.push({ title, content, time: new Date().toLocaleString() });

    // Xóa form
    document.getElementById('postTitle').value = '';
    document.getElementById('postContent').value = '';

    // Hiển thị lại bài đăng
    renderPosts();
}

// Hàm hiển thị bài đăng
function renderPosts() {
    const container = document.getElementById('postsContainer');
    container.innerHTML = ''; // Xóa nội dung cũ

    posts.slice().reverse().forEach(post => { // Hiển thị mới nhất lên trên
        const postDiv = document.createElement('div');
        postDiv.className = 'post';
        postDiv.innerHTML = `
            <h4>${post.title}</h4>
            <p>${post.content}</p>
            <small>Đăng lúc: ${post.time}</small>
            <hr>
        `;
        container.appendChild(postDiv);
    });
}
