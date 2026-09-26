#!/usr/bin/env python3
"""Generate editable draw.io SVGs using only Python's standard library.

Run from any directory. --check verifies that the generated files are current.
The API descriptions were reviewed against controller-runtime v0.25.1 and
client-go v0.37.1; update the descriptions as well as versions on migration.
"""
import argparse
from html import escape
from pathlib import Path
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[1] / 'contents/kubernetes-operator'
RUNTIME = 'controller-runtime v0.25.1'
CLIENT = 'client-go v0.37.1'
# Each card has a title and up to three explanatory lines. Edges describe
# dependency, event or call direction; their meaning is stated in each diagram.
DIAGRAMS = []


def diagram(path, title, version, cards, edges, note):
    DIAGRAMS.append((path, title, version, cards, edges, note))


def flow(path, title, version, cards, note):
    # Six steps, read left-to-right, then right-to-left on the lower row.
    diagram(path, title, version, [cards[i] for i in [0, 1, 2, 5, 4, 3]],
            [(0, 1), (1, 2), (2, 4), (4, 5), (4, 3), (5, 3)], note)


flow('controller.drawio.svg', 'The controller loop', CLIENT + ' / ' + RUNTIME, [
    ('Kubernetes API', 'Desired and observed state', 'List / Watch or WatchList'),
    ('Informer', 'Maintain a local cache', 'Deliver object notifications'),
    ('Event handler', 'Map an object to a key', 'Usually namespace / name'),
    ('Workqueue', 'Coalesce duplicate keys', 'Track work and retries'),
    ('Reconcile / sync', 'Read current state', 'Compute the required changes'),
    ('API client', 'Write only when necessary', 'Writes may trigger new events'),
], 'Reconciliation is state-based. A notification is not a durable event history.')
flow('controller-detailed.drawio.svg', 'From watch data to reconciliation', CLIENT, [
    ('Reflector', 'List + Watch / WatchList', 'Reconnect and refresh the Store'),
    ('Informer queue', 'Buffer object changes', 'Queue implementation may vary'),
    ('Indexer + notifications', 'Update the local indexed Store', 'Notify registered handlers'),
    ('Controller workqueue', 'Handlers enqueue object keys', 'Separate from the informer queue'),
    ('Worker + Lister', 'Get a key; read cached objects', 'DeepCopy before modification'),
    ('Clientset', 'Write changes to the API', 'Done / Forget or retry the key'),
], 'Informer queues update cached objects; controller workqueues schedule reconciliation keys.')
flow('client-go/diagram.drawio.svg', 'client-go controller data flow', CLIENT, [
    ('Clientset + ListWatch', 'Connect to Kubernetes API', 'Context-aware List / Watch'),
    ('Reflector', 'Populate and refresh the Store', 'Recover from watch interruptions'),
    ('SharedIndexInformer', 'Process queued changes', 'Update Indexer; notify handlers'),
    ('ResourceEventHandler', 'Extract namespace / name', 'Handle deletion tombstones'),
    ('Typed workqueue', 'Get key; call sync logic', 'Done for every acquired item'),
    ('Lister + API client', 'Read cache; write through client', 'Forget on success; retry errors'),
], 'Treat cached objects as read-only. Resync notifications do not imply a fresh API List.')
flow('client-go/clientset/clientset-simple.drawio.svg', 'Clientset: a typed API request', CLIENT, [
    ('Kubeconfig', 'Load rest.Config', 'Check configuration errors'),
    ('kubernetes.NewForConfig', 'Create the Clientset', 'Check client construction errors'),
    ('CoreV1() / AppsV1()', 'Select API group and version', 'For built-in resource types'),
    ('Pods(namespace)', 'Select resource and scope', 'Empty namespace: all for List'),
    ('List(ctx, options)', 'Send request through REST client', 'Set a timeout or cancellation'),
    ('PodList + error', 'Check error before using Items', 'Display namespace / name'),
], 'Custom resources require a generated client, dynamic client or controller-runtime Client.')
diagram('client-go/clientset/clientset.drawio.svg', 'Clientset interfaces and resource clients', CLIENT, [
    ('kubernetes.Interface', 'Typed group accessors', 'Discovery() accessor'),
    ('CoreV1Interface', 'Pods(namespace)', 'ConfigMaps(namespace), ...'),
    ('PodInterface', 'Get / List / Watch', 'Create / Update / Patch / Delete'),
    ('DiscoveryInterface', 'Discover groups and resources', 'Does not read Pod objects'),
    ('AppsV1Interface', 'Deployments(namespace)', 'ReplicaSets(namespace), ...'),
    ('DeploymentInterface', 'Typed Deployment operations', 'Context + options + error'),
], [(0, 1), (1, 2), (0, 3), (0, 4), (4, 5)],
    'Group accessors select clients. Resource accessors bind a namespace; methods send API requests.')
flow('client-go/informer/informer-factory.drawio.svg', 'SharedInformerFactory lifecycle', CLIENT, [
    ('NewSharedInformerFactory', 'Provide clientset + resync period', 'One shared informer per type'),
    ('Core().V1().Pods()', 'Obtain a typed informer wrapper', 'No running watch yet'),
    ('Informer() + handler', 'InformerFor creates / reuses it', 'Register handler; check error'),
    ('factory.Start(ctx.Done())', 'Start registered informers', 'Run watches in the background'),
    ('WaitForCacheSync', 'Wait for the initial cache state', 'Then use the typed Lister'),
    ('Cancel + Shutdown()', 'Cancel context to stop watches', 'Shutdown waits for goroutines'),
], 'Informer HasSynced covers cache initialization; a handler registration has its own HasSynced.')
flow('client-go/informer/informer.drawio.svg', 'Inside a SharedIndexInformer', CLIENT, [
    ('ListerWatcher + Reflector', 'Read initial state and changes', 'List + Watch or WatchList'),
    ('Internal queue', 'Buffer changes from Reflector', 'Chosen by implementation / flags'),
    ('Process changes', 'Apply Add / Update / Delete', 'Maintain the Indexer'),
    ('sharedProcessor', 'Distribute notifications', 'To registered listeners'),
    ('ResourceEventHandler', 'OnAdd / OnUpdate / OnDelete', 'Keep callbacks lightweight'),
    ('Controller workqueue', 'Enqueue keys for a worker', 'Worker reads via a Lister'),
], 'Lister and handlers share cached objects. DeepCopy before editing; handle deletion tombstones.')
flow('client-go/reflector/diagram.drawio.svg', 'Reflector: keep a Store up to date', CLIENT, [
    ('RunWithContext(ctx)', 'Start the reflector loop', 'Stop when context is canceled'),
    ('Initial state', 'ListWithContext or WatchList', 'Obtain a resourceVersion'),
    ('Store.Replace', 'Establish the initial snapshot', 'Store may be an informer queue'),
    ('Watch changes', 'Continue from resourceVersion', 'Process events from the API'),
    ('Store operations', 'Add / Update / Delete', 'Reflect changes downstream'),
    ('Reconnect / refresh', 'Recover from closed watches', 'Relist when required'),
], 'resourceVersion is opaque. Resync of known objects is different from relisting the API.')
flow('client-go/deltafifo/deltafifo.drawio.svg', 'DeltaFIFO: accumulate changes by key', CLIENT, [
    ('Producer', 'Add / Update / Delete', 'Replace / Resync also supported'),
    ('Key function', 'Map object to namespace / name', 'One queued entry per key'),
    ('items + queue', 'items: key to Deltas', 'queue: ordered keys'),
    ('Pop(process)', 'Wait for an available key', 'Return accumulated Deltas'),
    ('process(obj, initial)', 'initial: isInInitialList', 'Apply deltas to the consumer'),
    ('Consumer / Close()', 'Real informer updates Indexer', 'Close unblocks waiting Pop'),
], 'KnownObjects supplies existing state; it does not update the Indexer automatically. Sample only logs.')
diagram('client-go/indexer/indexer.drawio.svg', 'Indexer: objects, keys and secondary indexes', CLIENT, [
    ('Object', 'Example: Pod default/nginx', 'Input to key and index functions'),
    ('KeyFunc', 'MetaNamespaceKeyFunc', 'Produces default/nginx'),
    ('items', 'key -> object', 'GetByKey retrieves the object'),
    ('IndexFunc', 'Example: namespace index', 'Produces ["default"]'),
    ('Secondary index', 'index name + indexed value', 'Maps value to object keys'),
    ('ByIndex / Index', 'Resolve matching keys', 'Return stored objects'),
], [(0, 1), (1, 2), (0, 3), (3, 4), (4, 5), (2, 5)],
    'An Indexer extends Store. Results share stored objects; an index is not a copy of those objects.')
flow('controller-runtime/diagram.drawio.svg', 'controller-runtime: setup and reconciliation', RUNTIME, [
    ('Manager + Cluster', 'Create Client, Cache and Scheme', 'Pass dependencies explicitly'),
    ('Builder', 'For / Owns / Watches', 'Complete(reconciler)'),
    ('Controller + Source', 'Manager starts the controller', 'Start sources; wait for sync'),
    ('Handler + queue', 'Map events to request keys', 'Coalesce and schedule requests'),
    ('Reconciler', 'Reconcile(ctx, request)', 'Read current state via Client'),
    ('Client + Kubernetes API', 'Cached reads by default', 'Write directly to the API'),
], 'Manager also manages webhook and cache lifecycles. Dependencies are passed during construction.')
flow('controller-runtime/builder/overview.drawio.svg', 'Builder: construct and register a controller', RUNTIME, [
    ('Manager', 'Provide shared dependencies', 'GetClient / GetCache / GetScheme'),
    ('NewControllerManagedBy', 'Create a Builder for the Manager', 'Configure controller options'),
    ('For / Owns / Watches', 'Select objects and handlers', 'Set predicates as needed'),
    ('Complete(reconciler)', 'Pass initialized Reconciler', 'Build returns Controller + error'),
    ('Controller + sources', 'Register Controller with Manager', 'Watch the configured sources'),
    ('Manager.Start(ctx)', 'Start cache and controllers', 'Reconcile after source sync'),
], 'Complete returns an error; check it. A Client field is not automatically populated.')
diagram('controller-runtime/builder/for-owns-watches.drawio.svg', 'For, Owns and Watches choose request keys', RUNTIME, [
    ('For(primary)', 'Watch the primary object type', 'Use its namespace / name'),
    ('Owns(child)', 'Watch child objects', 'Map controller owner to its key'),
    ('Watches(object, handler)', 'Watch another resource type', 'Use a custom mapping handler'),
    ('source.Kind(...)', 'Cache + object + handler', 'Predicates supplied at construction'),
    ('Controller.Watch(src)', 'Register the complete Source', 'Start sources with typed queue'),
    ('Reconcile(ctx, request)', 'Request identifies target state', 'Event type is not in Request'),
], [(0, 3), (1, 3), (2, 3), (3, 4), (4, 5)],
    'All three Builder methods create sources; WatchesRawSource accepts an already constructed Source.')
diagram('controller-runtime/builder/for-owns-example.drawio.svg', 'ReplicaSet example: which object is reconciled?', RUNTIME, [
    ('ReplicaSet example', 'For(&appsv1.ReplicaSet{})', 'A change targets this ReplicaSet'),
    ('Owned Pod', 'Owns(&corev1.Pod{})', 'controller owner points to example'),
    ('Unrelated Pod', 'No matching controller owner', 'No request for this ReplicaSet'),
    ('Request', 'Namespace: default', 'Name: example'),
    ('ReplicaSetReconciler', 'Read ReplicaSet + matching Pods', 'Count only owned Pod UIDs'),
    ('pod-count label', 'Patch only when value changes', 'NotFound: finish successfully'),
], [(0, 3), (1, 3), (3, 4), (4, 5)],
    'Owns observes ownerReferences; it does not create them or create child Pods.')
flow('controller-runtime/builder/manager-perspective.drawio.svg', 'How a Builder controller is started', RUNTIME, [
    ('Builder.Complete(r)', 'Build Controller with Reconciler', 'Register configured sources'),
    ('Manager.Add(controller)', 'Classify the Runnable', 'NeedLeaderElection is consulted'),
    ('LeaderElection group', 'Default group for Controllers', 'Not the generic Others group'),
    ('Manager.Start(ctx)', 'Start infrastructure and caches', 'Manage leadership when enabled'),
    ('Controller.Start(ctx)', 'Start Source; wait for sync', 'Run reconciliation workers'),
    ('Context cancellation', 'Stop workers and sources', 'Manager coordinates shutdown'),
], 'With leader election disabled, controllers still run. Non-leader runnables can run on every replica.')
diagram('controller-runtime/manager/diagram.drawio.svg', 'Manager owns shared dependencies and lifecycles', RUNTIME, [
    ('Manager', 'Embeds the Cluster interface', 'Add(Runnable); Start(ctx)'),
    ('Cluster dependencies', 'Client / Cache / APIReader', 'Scheme / RESTMapper / Config'),
    ('Reconciler construction', 'Client: mgr.GetClient()', 'Pass other dependencies explicitly'),
    ('Infrastructure runnables', 'HTTPServers / Webhooks / Caches', 'Others and Warmup groups'),
    ('LeaderElection runnables', 'Controllers by default', 'Run according to leadership'),
    ('Shutdown', 'Cancel context; stop runnables', 'Use SetupSignalHandler()'),
], [(0, 1), (1, 2), (0, 3), (3, 4), (4, 5)],
    'Lifecycle view, not an exhaustive startup sequence. Check Add, Complete and Start errors.')
flow('controller-runtime/controller/diagram.drawio.svg', 'Controller workers and retry behavior', RUNTIME, [
    ('Watch(src)', 'Register a constructed Source', 'Handler belongs to the Source'),
    ('Start(ctx)', 'Start sources with the queue', 'Wait for SyncingSource readiness'),
    ('Typed request queue', 'Handlers enqueue request keys', 'Workers call Get()'),
    ('Reconcile(ctx, request)', 'Read and adjust current state', 'Return Result and error'),
    ('Choose next action', 'Success: Forget; delay: AddAfter', 'Error: rate-limited retry'),
    ('Done(request)', 'Finish processing this key', 'Process newly queued work'),
], 'RequeueAfter schedules an explicit delay. TerminalError does not cause an automatic error retry.')
diagram('controller-runtime/reconciler/diagram.drawio.svg', 'Reconciler inputs and outcomes', RUNTIME, [
    ('Controller worker', 'Get a request from the queue', 'Invoke the configured Reconciler'),
    ('Request', 'types.NamespacedName', 'No old object or event type'),
    ('Reconcile(ctx, request)', 'Read state and adjust it', 'Make repeated calls safe'),
    ('Success', 'Result{}, nil', 'Forget rate-limiter state'),
    ('Scheduled retry', 'Result{RequeueAfter: duration}', 'Return nil error'),
    ('Error', 'Normal error: rate-limited retry', 'TerminalError: no error retry'),
], [(0, 1), (1, 2), (2, 3), (2, 4), (2, 5)],
    'Controller calls Done for the request. New events can trigger reconciliation after any outcome.')
diagram('controller-runtime/client/diagram.drawio.svg', 'Manager Client: read and write paths', RUNTIME, [
    ('Reconciler', 'Receives mgr.GetClient()', 'Uses Get / List and writes'),
    ('client.Client', 'client.Options.Cache.Reader', 'CacheOptions controls bypass'),
    ('Cache reader', 'Get / List for cached objects', 'Informer-backed, eventually current'),
    ('mgr.GetAPIReader()', 'Direct API Get / List', 'Explicit uncached reads'),
    ('Direct API operations', 'Writes and Status operations', 'Uncached Get / List'),
    ('Kubernetes API server', 'Receives direct reads and writes', 'Watches later refresh the cache'),
], [(0, 1), (1, 2), (0, 3), (1, 4), (3, 5), (4, 5)],
    'A standalone client.New without Cache reads directly from the API. Writes are not immediately cached.')
flow('controller-runtime/cache/diagram.drawio.svg', 'Standalone Cache: registration, sync and reads', RUNTIME, [
    ('cache.New(config, options)', 'Provide Scheme as needed', 'HTTP client + mapper defaulted'),
    ('GetInformer(ctx, object)', 'Register a watched resource', 'Before starting this example'),
    ('AddEventHandler', 'Register callbacks', 'Check returned error'),
    ('Start(ctx)', 'Run cache in another goroutine', 'Propagate startup errors'),
    ('WaitForCacheSync(ctx)', 'Wait for the initial state', 'Do not ignore false'),
    ('Get / List', 'Read from the cache', 'Cancel context to stop'),
], 'Do not call Get on an empty object to initialize a cache. Read only after startup and synchronization.')
diagram('controller-runtime/cache/diagram-2.drawio.svg', 'Informer Cache: implementation layers', RUNTIME, [
    ('cache.Cache interface', 'client.Reader + Informers', 'Get / List / GetInformer / Start'),
    ('informerCache', 'Uses internal.Informers', 'Resolve type and obtain a reader'),
    ('internal.Informers tracker', 'Structured / Unstructured / Metadata', 'Each maps GVK to *internal.Cache'),
    ('client-go Indexer', 'Local indexed object Store', 'Populated by SharedIndexInformer'),
    ('internal.Cache', 'Informer: SharedIndexInformer', 'Reader: CacheReader'),
    ('CacheReader', 'Uses the informer Indexer', 'Get / List; copies by default'),
], [(0, 1), (1, 2), (2, 4), (4, 5), (4, 3), (5, 3)],
    'Single informer-cache view. Namespace and per-object options can add multi-namespace / delegating caches.')
diagram('controller-runtime/source/diagram.drawio.svg', 'Source API: dependencies are explicit', RUNTIME, [
    ('TypedSource[Request]', 'Start(ctx, typedQueue)', 'Source = TypedSource[Request]'),
    ('TypedSyncingSource', 'TypedSource + WaitForSync(ctx)', 'Controller waits before workers'),
    ('Kind / TypedKind', 'Construct with cache and object', 'Also provide handler + predicates'),
    ('Channel / TypedChannel', 'GenericEvents from a channel', 'Handler supplied at construction'),
    ('Informer / TypedInformer', 'Existing informer + handler', 'Lower-level source adapter'),
    ('Controller', 'Watch(src)', 'Start source; wait when supported'),
], [(1, 0), (2, 1), (3, 0), (4, 0), (5, 0)],
    'For Source aliases, Request is reconcile.Request. Arrows point to an implemented interface or the interface being called.')
diagram('controller-runtime/source/dataflow.drawio.svg', 'Sources feed handlers and the request queue', RUNTIME, [
    ('Kubernetes events', 'Kind uses Cache + Informer', 'Informer source accepts one directly'),
    ('Predicates', 'Filter incoming events', 'Configured on the Source'),
    ('TypedEventHandler', 'Map event objects to requests', 'For self, owner or custom targets'),
    ('External events', 'User code sends GenericEvents', 'Channel source receives them'),
    ('Typed request queue', 'Store comparable request keys', 'ShutDown unblocks Get()'),
    ('Reconciler / worker', 'Get key; process current state', 'Forget on success; always Done'),
], [(0, 1), (1, 2), (3, 1), (2, 4), (4, 5)],
    'A source starts with Start(ctx, queue). It already holds its handler and predicates.')

SVG_NS = 'http://www.w3.org/2000/svg'
WIDTH, HEIGHT = 1240, 610
CARD_W, CARD_H = 330, 128
POSITIONS = [(45 + col * 410, 124 + row * 236) for row in range(2) for col in range(3)]


def route(a, b, offset=0):
    """Route arrows through the wide gutters, keeping them out of cards."""
    x1, y1 = POSITIONS[a]
    x2, y2 = POSITIONS[b]
    if y1 == y2:
        if abs(a - b) == 1:
            return [(x1 + (CARD_W if x2 > x1 else 0), y1 + CARD_H / 2 + offset),
                    (x2 + (0 if x2 > x1 else CARD_W), y2 + CARD_H / 2 + offset)]
        # A nonadjacent horizontal connection goes through the row gutter.
        y = y1 + CARD_H + 42
        return [(x1 + CARD_W / 2 + offset, y1 + CARD_H), (x1 + CARD_W / 2 + offset, y),
                (x2 + CARD_W / 2 + offset, y), (x2 + CARD_W / 2 + offset, y2 + CARD_H)]
    if x1 == x2:
        return [(x1 + CARD_W / 2 + offset, y1 + (CARD_H if y2 > y1 else 0)),
                (x2 + CARD_W / 2 + offset, y2 + (0 if y2 > y1 else CARD_H))]
    middle = 290 + (a * 7 + b * 3) % 45
    return [(x1 + CARD_W / 2 + offset, y1 + (CARD_H if y2 > y1 else 0)),
            (x1 + CARD_W / 2 + offset, middle), (x2 + CARD_W / 2 + offset, middle),
            (x2 + CARD_W / 2 + offset, y2 + (0 if y2 > y1 else CARD_H))]


def generate(title, version, cards, edges, note):
    mxfile = ET.Element('mxfile', host='app.diagrams.net', type='device')
    page = ET.SubElement(mxfile, 'diagram', id='main', name=title)
    model = ET.SubElement(page, 'mxGraphModel', dx=str(WIDTH), dy=str(HEIGHT), grid='1', gridSize='10',
                          page='1', pageScale='1', pageWidth=str(WIDTH), pageHeight=str(HEIGHT))
    root = ET.SubElement(model, 'root')
    ET.SubElement(root, 'mxCell', id='0')
    ET.SubElement(root, 'mxCell', id='1', parent='0')

    def cell(identifier, value, x, y, w, h, style, parent='1'):
        c = ET.SubElement(root, 'mxCell', id=identifier, value=value, style=style, vertex='1', parent=parent)
        ET.SubElement(c, 'mxGeometry', x=str(x), y=str(y), width=str(w), height=str(h), attrib={'as': 'geometry'})
        return c

    text_style = 'text;html=0;align=left;verticalAlign=middle;whiteSpace=wrap;fontFamily=Arial;'
    cell('title', title, 45, 25, 1150, 40, text_style + 'fontSize=28;fontStyle=1;fontColor=#172B4D;')
    cell('version', version + '  |  Read arrows in their direction', 45, 75, 1150, 24,
         text_style + 'fontSize=16;fontColor=#52677F;')
    cell('note', note, 45, 535, 1150, 48, text_style + 'fontSize=15;fontColor=#52677F;')
    parts = [f'<svg xmlns="{SVG_NS}" width="{WIDTH}" height="{HEIGHT}" viewBox="0 0 {WIDTH} {HEIGHT}" role="img" aria-labelledby="title desc">',
             f'<title id="title">{escape(title)}</title>', f'<desc id="desc">{escape(note)}</desc>',
             '<defs><marker id="arrow" markerWidth="9" markerHeight="9" refX="8" refY="4.5" orient="auto" markerUnits="userSpaceOnUse"><path d="M0,0 L9,4.5 L0,9 Z" fill="#52677f"/></marker></defs>',
             '<rect width="1240" height="610" fill="#ffffff"/>',
             f'<text x="45" y="57" font-family="Arial, sans-serif" font-size="28" font-weight="700" fill="#172b4d">{escape(title)}</text>',
             f'<text x="45" y="94" font-family="Arial, sans-serif" font-size="16" fill="#52677f">{escape(version)}  |  Read arrows in their direction</text>']
    for i, (a, b) in enumerate(edges):
        # Distinct ports prevent diagonal connections from merging into
        # unrelated vertical arrows in the same column.
        complex_routes = any(u // 3 != v // 3 and u % 3 != v % 3 for u, v in edges)
        offset = (i - (len(edges) - 1) / 2) * 14 if complex_routes else 0
        points = route(a, b, offset)
        parts.append('<path d="' + ' '.join(('M' if j == 0 else 'L') + f'{x:g},{y:g}' for j, (x, y) in enumerate(points)) + '" fill="none" stroke="#52677f" stroke-width="2" marker-end="url(#arrow)"/>')
        sx, sy = POSITIONS[a]
        tx, ty = POSITIONS[b]
        anchors = (f'exitX={(points[0][0]-sx)/CARD_W};exitY={(points[0][1]-sy)/CARD_H};'
                   f'entryX={(points[-1][0]-tx)/CARD_W};entryY={(points[-1][1]-ty)/CARD_H};')
        c = ET.SubElement(root, 'mxCell', id=f'e{i}', edge='1', parent='1', source=f'n{a}', target=f'n{b}',
                          style=anchors + 'edgeStyle=none;html=0;endArrow=block;endFill=1;strokeColor=#52677F;strokeWidth=2;')
        geom = ET.SubElement(c, 'mxGeometry', relative='1', attrib={'as': 'geometry'})
        if len(points) > 2:
            array = ET.SubElement(geom, 'Array', attrib={'as': 'points'})
            for x, y in points[1:-1]:
                ET.SubElement(array, 'mxPoint', x=str(x), y=str(y))
    for i, (heading, *body) in enumerate(cards):
        x, y = POSITIONS[i]
        cell(f'n{i}', '', x, y, CARD_W, CARD_H,
             'rounded=1;arcSize=12;html=0;fillColor=#F3F7FC;strokeColor=#B7C9E2;strokeWidth=1;')
        cell(f'n{i}title', heading, 18, 14, CARD_W - 36, 30,
             text_style + 'fontSize=18;fontStyle=1;fontColor=#17375E;', f'n{i}')
        cell(f'n{i}body', '\n'.join(body), 18, 51, CARD_W - 36, 63,
             text_style + 'fontSize=16;fontColor=#334E68;', f'n{i}')
        parts.append(f'<rect x="{x}" y="{y}" width="{CARD_W}" height="{CARD_H}" rx="12" fill="#f3f7fc" stroke="#b7c9e2"/>')
        parts.append(f'<text x="{x+18}" y="{y+35}" font-family="Arial, sans-serif" font-size="18" font-weight="700" fill="#17375e">{escape(heading)}</text>')
        for j, line in enumerate(body):
            parts.append(f'<text x="{x+18}" y="{y+72+j*25}" font-family="Arial, sans-serif" font-size="16" fill="#334e68">{escape(line)}</text>')
    # Two lines keep the explanatory note readable at normal README widths.
    import textwrap
    for i, line in enumerate(textwrap.wrap(note, width=132)):
        parts.append(f'<text x="45" y="{555+i*23}" font-family="Arial, sans-serif" font-size="15" fill="#52677f">{escape(line)}</text>')
    parts.append('</svg>')
    svg = ET.fromstring('\n'.join(parts))
    svg.set('content', ET.tostring(mxfile, encoding='unicode'))
    ET.register_namespace('', SVG_NS)
    return ET.tostring(svg, encoding='unicode') + '\n'


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--check', action='store_true')
    args = parser.parse_args()
    stale = []
    for path, title, version, cards, edges, note in DIAGRAMS:
        target = ROOT / path
        content = generate(title, version, cards, edges, note)
        if args.check:
            if not target.exists() or target.read_text() != content:
                stale.append(path)
        else:
            target.write_text(content)
    if stale:
        parser.exit(1, 'Stale diagrams:\n' + '\n'.join(stale) + '\n')
    print(f'{"Checked" if args.check else "Generated"} {len(DIAGRAMS)} editable SVG diagrams.')


if __name__ == '__main__':
    main()
