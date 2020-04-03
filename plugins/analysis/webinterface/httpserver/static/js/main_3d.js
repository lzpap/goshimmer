const ANALYSIS_SERVER_URL = "127.0.0.1" + "/datastream";
const NODE_ID_LENGTH = 64;

class Frontend {
    constructor(app) {
        this.app = app;
        this.activeNode = '';
        this.searchTerm = '';
    }

    // Show search field
    showSearchField() {
        document.getElementById('searchWrapper').style.cssText = "display:block;"
    }

    // Show nodes
    displayNodesOnline() {
        // this line might be a performance killer!

        let html = [];
        for(let n of this.app.ds.gData.nodes) {
            if(n.id == this.activeNode) {
                html.push('<span class="n active">' + n.id + "</span>");
            } else {
                html.push('<span class="n">' + n.id + "</span>");
            }
        }
        document.getElementById("nodesOnline").innerHTML = html.join("");
    }
}

class Datastructure {
    constructor(app) {
        this.app = app;
        this.gData = {
            nodes: [],
            links: []
        }; 
    }

    addNode(idA) {
        this.gData.nodes.push({ id: idA})
        this.app.graph.setGraphData(this.gData);
    }

    removeNode(idA) {
        this.gData.links = this.gData.links.filter(l => l.source.id !== idA && l.target.id !== idA); // Remove links attached to node
        this.gData.nodes = this.gData.nodes.filter(n => n.id !== idA); // Remove node
        this.app.graph.setGraphData(this.gData);
    }

    connectNodes(idA, idB) {
        this.gData.links.push({ source: idA, target: idB })
        this.app.graph.setGraphData(this.gData);

    }
    disconnectNodes(idA, idB) {
        this.gData.links = this.gData.links.filter(l => (l.source.id !== idA && l.target.id !== idB));
        this.gData.links = this.gData.links.filter(l => (l.source.id !== idB && l.target.id !== idA));
        this.app.graph.setGraphData(this.gData);
    }
}

class Graph{
    constructor(app) {
        this.app = app
        this.graph = ForceGraph3D()
            (document.getElementById('graphc'))
            .graphData(app.ds.gData)
            .linkDirectionalParticles(5)
            .linkDirectionalParticleSpeed(0.03);
    }

    setGraphData(gData) {
        this.graph.graphData(gData)
    }

    getGraphData() {
        return this.graph.graphData()
    }
}

class Application {
    constructor(url) {
        this.url = url;
        this.ds = new Datastructure(this)
        this.graph = new Graph(this)
        this.frontend = new Frontend(this)
    }

    run() {
        this.initWebsocket();
        this.showOnlineNodes();
    }

    initWebsocket() {
        this.socket = new WebSocket(
            ((window.location.protocol === "https:") ? "wss://" : "ws://") + this.url
        );
    
        this.socket.onopen = () => {
            setInterval(() => {
                this.socket.send("_");
            }, 1000);
        };
    
        this.socket.onerror = (e) => {
            console.error("WebSocket error observed", e);
          };
    
        this.socket.onmessage = (e) => {
            let type = e.data[0];
            let data = e.data.substr(1);
            let idA = data.substr(0, NODE_ID_LENGTH);
            let idB;


            switch (type) {
                case "_":
                    //do nothing - its just a ping
                    break;
    
                case "A":
                    console.log("addNode event:", idA);
                    // filter out empty ids
                    if(idA.length == NODE_ID_LENGTH) {
                        this.ds.addNode(idA);
                    }
                    break;

                case "a":
                    console.log("removeNode event:", idA);
                    this.ds.removeNode(idA);
                    break;
    
                case "C":
                    idB = data.substr(NODE_ID_LENGTH, NODE_ID_LENGTH);
                    console.log("connectNodes event:", idA, " - ", idB);
                    this.ds.connectNodes(idA, idB);
                    break;
    
                case "c":
                    idB = data.substr(NODE_ID_LENGTH, NODE_ID_LENGTH);
                    console.log("disconnectNodes event:", idA, " - ", idB);
                    this.ds.disconnectNodes( idA, idB);
                    break;
            }
        }
    }

    showOnlineNodes() {
        setInterval(() => { 
            if(this.frontend.searchTerm.length > 0) {
                return;
            }
            this.frontend.displayNodesOnline() 
        }, 300);
    }
}

let app;
window.onload = () => {
    app = new Application(ANALYSIS_SERVER_URL);
    app.run()
}