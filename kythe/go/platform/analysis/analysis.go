/*
 * Copyright 2014 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package analysis defines interfaces used to locate and analyze compilation
// units.
//
// Implementations of the Analyzer interface can be plugged in to analyze
// targets from compilation data stored in a variety of formats.
//
// Implementations of the Fetcher interface permit an Analyzer to consume files
// from index files, local files, or other data sources.
//
// Implementations of the Sink interface provide a means for an Analyzer to
// emit output that can be captured as artifacts of the analysis process.
//
// Implementations of the Runner interface invoke an Analyzer on a collection
// of compilation units.
package analysis

import (
    "errors"

    "code.google.com/p/goprotobuf/proto"

    apb "kythe/proto/analysis_proto"

)

// An Analyzer provides the ability to perform arbitrary analysis on a single
// compilation unit, possibly emitting analysis artifacts as output via the Sink.
type Analyzer interface {
    // Analyze performs analysis of a single compilation unit.  The analyzer
    // can retrieve file data using the Fetcher, and can emit analysis
    // artifacts via the Sink.  Returns an error if analysis did not succeed.
    Analyze(*apb.AnalysisRequest, Fetcher, Sink) error
}

// A Sink captures data artifacts generated by an Analyzer for storage or
// transmission.
type Sink interface {
    // WriteBytes emits an arbitrary slice of bytes.  Each write to a sink is
    // treated as a single opaque artifact, typically a serialized proto.
    WriteBytes([]byte) error
}

// WriteMessage is a helper function that marshals an arbitrary protocol buffer
// message and writes its marshalled form to sink.
func WriteMessage(msg proto.Message, sink Sink) error {
    // Just as a safety measure, treat a nil sink as a no-op write.
    if sink == nil {
        return nil
    }

    data, err := proto.Marshal(msg)
    if err != nil {
        return err
    }
    return sink.WriteBytes(data)
}

// A Fetcher provides the ability to fetch the contents of specified files.
type Fetcher interface {
    // Fetch retrieves the contents of a single file.  At least one of path and
    // digest must be provided; both are preferred.  The implementation decides
    // what to do if only one is given.
    Fetch(path, digest string) ([]byte, error)
}

// ErrNotFound is returned by Fetch when the specified file was not found.
var ErrNotFound = errors.New("file not found")

// A Runner invokes an Analyzer on a collection of compilation units.
type Runner interface {
    // RunAnalysis runs analyzer on each compilation known to the runner.
    //
    // If an error occurs in the running process, such as inability to read
    // compilation unit data, it is returned.
    //
    // Errors in analysis are not returned from RunAnalysis: If report != nil,
    // it is called for each error returned by analyzer along with the analysis
    // request being processed when the error occurred.
    RunAnalysis(analyzer Analyzer, report func(error, *apb.AnalysisRequest)) error
}