#pragma once

#ifndef LOG_LEVEL
#define LOG_LEVEL 3
#warning "No log level set, defaulting to 3 (INFO)"
#endif

#if LOG_LEVEL >= 1
#define ERROR(...) printk("ERROR: " __VA_ARGS__)
#else
#define ERROR(...) 
#endif

#if LOG_LEVEL >= 2
#define WARN(...) printk("WARN : " __VA_ARGS__)
#else
#define WARN(...) 
#endif

#if LOG_LEVEL >= 3
#define INFO(...) printk("INFO : " __VA_ARGS__)
#else
#define INFO(...) 
#endif

#if LOG_LEVEL >= 4
#define DEBUG(...) printk("DEBUG: " __VA_ARGS__)
#else
#define DEBUG(...) 
#endif

#if LOG_LEVEL >= 5
#define TRACE(...) printk("TRACE: " __VA_ARGS__)
#else
#define TRACE(...) 
#endif
