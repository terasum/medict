
import { contentTracing } from 'electron';
 // Setting custom category to get stack traces
  const defaultTraceCategories: Readonly<Array<string>> = [
    '-*', 'devtools.timeline', 'disabled-by-default-devtools.timeline',
    'disabled-by-default-devtools.timeline.frame',
    'toplevel', 'blink.console',
    'disabled-by-default-devtools.timeline.stack',
    // 'disabled-by-default-v8.cpu_profile', 
    // 'disabled-by-default-v8.cpu_profiler',
    // 'disabled-by-default-v8.cpu_profiler.hires'
  ];

  const traceOptions = {
    categoryFilter: defaultTraceCategories.join(','),
    traceOptions: 'record-until-full',
    options: 'sampling-frequency=10000'
  };


  // contentTracing.startRecording(traceOptions)

