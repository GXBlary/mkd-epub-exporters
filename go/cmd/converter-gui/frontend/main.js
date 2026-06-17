import { SelectFiles, SelectFolder, OpenFile, Convert } from './wailsjs/go/main/App.js';
import { EventsOn } from './wailsjs/runtime/runtime.js';

// Application State
let queue = [];
let outputDir = "";

// DOM Elements
const dropZone = document.getElementById('dropZone');
const fileQueue = document.getElementById('fileQueue');
const queueCount = document.getElementById('queueCount');
const btnClearQueue = document.getElementById('btnClearQueue');
const btnBrowseOut = document.getElementById('btnBrowseOut');
const outputDirInput = document.getElementById('outputDir');
const btnConvert = document.getElementById('btnConvert');

// Progress DOM
const progressOverlay = document.getElementById('progressOverlay');
const progressStatus = document.getElementById('progressStatus');
const progressDetail = document.getElementById('progressDetail');
const progressBar = document.getElementById('progressBar');

// Results DOM
const resultsModal = document.getElementById('resultsModal');
const resSuccessCount = document.getElementById('resSuccessCount');
const resFailedCount = document.getElementById('resFailedCount');
const resultsList = document.getElementById('resultsList');
const btnCloseResults = document.getElementById('btnCloseResults');

// Drag and Drop Handlers
dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropZone.classList.add('active');
});

dropZone.addEventListener('dragleave', () => {
    dropZone.classList.remove('active');
});

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('active');
    
    const paths = [];
    if (e.dataTransfer.files) {
        for (let i = 0; i < e.dataTransfer.files.length; i++) {
            const file = e.dataTransfer.files[i];
            if (file.path) {
                paths.push(file.path);
            }
        }
    }
    
    if (paths.length > 0) {
        addToQueue(paths);
    }
});

// Click on DropZone to browse
dropZone.addEventListener('click', async () => {
    try {
        const selected = await SelectFiles();
        if (selected && selected.length > 0) {
            addToQueue(selected);
        }
    } catch (err) {
        console.error("Failed to select files:", err);
    }
});

// Browse Output Folder
btnBrowseOut.addEventListener('click', async () => {
    try {
        const selected = await SelectFolder();
        if (selected) {
            outputDir = selected;
            outputDirInput.value = selected;
        }
    } catch (err) {
        console.error("Failed to select folder:", err);
    }
});

// Clear Queue
btnClearQueue.addEventListener('click', () => {
    queue = [];
    updateQueueUI();
});

// Queue UI Updater
function addToQueue(paths) {
    paths.forEach(p => {
        if (!queue.includes(p)) {
            queue.push(p);
        }
    });
    updateQueueUI();
}

function updateQueueUI() {
    fileQueue.innerHTML = '';
    
    if (queue.length === 0) {
        fileQueue.innerHTML = `
            <div class="empty-queue-message">
                Queue is empty. Drop files above to start.
            </div>
        `;
        btnClearQueue.style.display = 'none';
        queueCount.innerText = '0';
        return;
    }
    
    btnClearQueue.style.display = 'block';
    queueCount.innerText = queue.length;
    
    queue.forEach((p, index) => {
        const div = document.createElement('div');
        div.className = 'queue-item';
        
        const filename = p.split(/[\\/]/).pop();
        
        div.innerHTML = `
            <div class="item-info">
                <span class="item-name" title="${p}">${filename}</span>
                <span class="item-meta">${p}</span>
            </div>
            <button class="btn-remove-item" data-index="${index}">&times;</button>
        `;
        
        div.querySelector('.btn-remove-item').addEventListener('click', (e) => {
            const idx = parseInt(e.target.getAttribute('data-index'));
            queue.splice(idx, 1);
            updateQueueUI();
        });
        
        fileQueue.appendChild(div);
    });
}

// Convert action
btnConvert.addEventListener('click', async () => {
    if (queue.length === 0) {
        alert("Please add at least one file or folder to the queue.");
        return;
    }
    if (!outputDir) {
        alert("Please select an output folder.");
        return;
    }
    
    const exportMD = document.getElementById('exportMD').checked;
    const exportEpub = document.getElementById('exportEpub').checked;
    const embedImages = document.getElementById('embedImages').checked;
    
    if (!exportMD && !exportEpub) {
        alert("Please select at least one export format (Markdown or EPUB).");
        return;
    }
    
    const format = (exportMD && exportEpub) ? "both" : (exportMD ? "md" : "epub");
    
    // Reset progress
    progressOverlay.style.display = 'flex';
    progressStatus.innerText = "Starting conversion...";
    progressDetail.innerText = `File 0/${queue.length}`;
    progressBar.style.width = '0%';
    
    try {
        const results = await Convert(queue, outputDir, format, embedImages);
        showResults(results);
    } catch (err) {
        alert("An error occurred during conversion: " + err);
        progressOverlay.style.display = 'none';
    }
});

// Show Results Modal
function showResults(results) {
    progressOverlay.style.display = 'none';
    
    let successCount = 0;
    let failedCount = 0;
    
    resultsList.innerHTML = '';
    
    results.forEach(res => {
        if (res.Success) {
            successCount++;
            const div = document.createElement('div');
            div.className = 'result-item';
            div.innerText = res.DestPath;
            div.title = "Double-click to open file";
            
            // Double click opens the file via Go App backend
            div.addEventListener('dblclick', () => {
                OpenFile(res.DestPath);
            });
            resultsList.appendChild(div);
        } else {
            failedCount++;
        }
    });
    
    resSuccessCount.innerText = successCount;
    resFailedCount.innerText = failedCount;
    
    resultsModal.style.display = 'flex';
}

btnCloseResults.addEventListener('click', () => {
    resultsModal.style.display = 'none';
    queue = []; // Clear queue on successful conversion
    updateQueueUI();
});

// Listen to progress events from Go Backend
EventsOn("conversion-progress", (data) => {
    const { current, total, name } = data;
    const pct = Math.round((current / total) * 100);
    
    progressStatus.innerText = `Converting: ${name}`;
    progressDetail.innerText = `File ${current}/${total}`;
    progressBar.style.width = `${pct}%`;
});
